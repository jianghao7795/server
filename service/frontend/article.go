package frontend

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"net/url"
	global "server/model"
	"server/model/frontend"
	"server/model/system"
	"strconv"
	"strings"
	"time"

	// json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"

	frontendReq "server/model/frontend/request"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Article struct{}

func encodeArticleListCache(list []frontend.Article) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&list); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeArticleListCache(data []byte) ([]frontend.Article, error) {
	var list []frontend.Article
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func encodeTotalCache(total int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(total))
	return b
}

func decodeTotalCache(b []byte) int64 {
	if len(b) < 8 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(b))
}

func (s *Article) GetArticleList(info *frontendReq.ArticleSearch, c fiber.Ctx) (list []frontend.Article, total int64, err error) {
	cacheTime := global.CONFIG.Cache.Time
	var articleBlob []byte
	pageKey := "article-list-" + strconv.Itoa(info.Page)
	db := global.DB.Model(&frontend.Article{})
	err = db.Count(&total).Error
	if err != nil {
		return list, 0, errors.New("总数请求失败")
	}
	switch info.IsImportant {
	case 1:
		articleBlob, err = global.REDIS.Get(c.Context(), "article-list-home").Bytes()
	case 2:
		articleBlob, err = global.REDIS.Get(c.Context(), pageKey).Bytes()
	default:
		articleBlob, err = global.REDIS.Get(c.Context(), pageKey).Bytes()
	}

	if errors.Is(err, redis.Nil) {
		limit := info.PageSize
		offset := info.PageSize * (info.Page - 1)

		if info.IsImportant != 0 {
			db = db.Where("is_important = ?", info.IsImportant)
		}
		if info.Title != "" {
			db = db.Where("title like ?", strings.Join([]string{"%", info.Title, "%"}, ""))
		}
		err = db.Limit(limit).Offset(offset).Order("id desc").Preload("Tags").Preload("User").Find(&list).Error
		if err != nil {
			return list, 0, err
		}

		listBytes, encErr := encodeArticleListCache(list)
		if encErr != nil {
			return list, 0, encErr
		}
		totalBytes := encodeTotalCache(total)

		if info.IsImportant != 0 {
			err = global.REDIS.Set(c.Context(), "article-list-home", listBytes, time.Duration(cacheTime)*time.Second).Err()
			if err != nil {
				return list, 0, errors.New("redis 存储失败")
			}
		} else {
			err = global.REDIS.Set(c.Context(), "article-list-total", totalBytes, time.Duration(cacheTime)*time.Second).Err()
			if err != nil {
				return list, 0, errors.New("redis 存储失败")
			}
			err = global.REDIS.Set(c.Context(), pageKey, listBytes, time.Duration(cacheTime)*time.Second).Err()
			if err != nil {
				return list, 0, errors.New("redis 存储失败")
			}
		}
	} else if err != nil {
		return list, 0, err
	} else if len(articleBlob) > 0 {
		list, err = decodeArticleListCache(articleBlob)
		if err != nil {
			return list, 0, err
		}
		if info.IsImportant == 0 {
			totalBlob, terr := global.REDIS.Get(c.Context(), "article-list-total").Bytes()
			if terr == nil && len(totalBlob) >= 8 {
				total = decodeTotalCache(totalBlob)
			}
		}
		return list, total, nil
	}

	return list, total, err
}

func (s *Article) GetArticleDetail(articleId int, c fiber.Ctx) (articleDetail frontend.Article, err error) {
	reqIP := c.IP()
	var ipUser frontend.Ip
	t := time.Now()
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	db := global.DB.Model(&frontend.Article{})
	dbIp := global.DB.Model(&frontend.Ip{}).Where("ip = ? and article_id = ?", reqIP, articleId).Where("created_at > ?", startTime).First(&ipUser)
	if errors.Is(dbIp.Error, gorm.ErrRecordNotFound) {
		locals := c.Locals("frontend_user")
		if locals != nil {
			ipUser.UserID = locals.(system.SysUser).ID
		} else {
			ipUser.UserID = 0
		}
		ipUser.ArticleID = uint(articleId)
		ipUser.Ip = reqIP

		err = global.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&ipUser).Error; err != nil {
				return err
			}
			if err = tx.Model(&frontend.Article{}).Where("id = ?", articleId).Update("reading_quantity", gorm.Expr("reading_quantity + ?", 1)).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return articleDetail, err
		}
	}
	err = db.Where("id = ?", articleId).Preload("Tags").Preload("User").First(&articleDetail).Error
	return articleDetail, err

	// var cacheTime = global.CONFIG.Cache.Time
	// var articleDetailStr string
	// articleDetailStr, err = global.REDIS.Get(c, "article"+strconv.Itoa(articleId)).Result()
	// if err == redis.Nil {
	// 	db := global.DB.Model(&frontend.Article{})
	// 	err = db.Where("id = ?", articleId).Preload("Tags").First(&articleDetail).Error
	// 	if err != nil {
	// 		return articleDetail, err
	// 	}
	// 	articleString, _ := json.Marshal(articleDetail)
	// 	err := global.REDIS.Set(c, "article"+strconv.Itoa(articleId), articleString, time.Duration(cacheTime)*time.Second).Err()
	// 	if err != nil {
	// 		global.LOG.Error("Redis 存储失败!", zap.Error(err))
	// 	}
	// 	return articleDetail, err
	// } else if err != nil {
	// 	return
	// } else {
	// 	if articleDetailStr != "" {
	// 		err = json.Unmarshal([]byte(articleDetailStr), &articleDetail)
	// 		return articleDetail, err
	// 	}
	// }

	// return
}

func (s *Article) GetSearchArticle(info *frontendReq.ArticleSearch) (list []frontend.Article, err error) {
	db := global.DB.Model(&frontend.Article{})
	// Preload("Tags", func(dbg *gorm.DB) *gorm.DB {
	// 	return dbg.Where("name = ?", info.Value)
	// })
	// 	type User struct {
	//     gorm.Model
	//     Name      string
	//     Phone     string
	//     Languages []*Language `gorm:"many2many:user_languages;"`
	// }

	// // Language 一种语言属于多个用户，使用 `user_languages` 作为连接表
	// type Language struct {
	//     gorm.Model
	//     Name  string
	//     Users []*User `gorm:"many2many:user_languages;"`
	// }
	sortField := make(map[string]string)
	sortField["read"] = "reading_quantity"
	sortField["time"] = "created_at"
	text, err := url.QueryUnescape(info.Value)
	if err != nil {
		return nil, err
	}
	if info.Name == "tags" {
		// 多对多关联 Association
		var id uint
		err = global.DB.Model(&frontend.Tag{}).Select("id").Where("name = ?", text).First(&id).Error
		if err != nil {
			return nil, err
		}
		dbTag := &frontend.Tag{MODEL: global.MODEL{ID: id}}
		if info.Sort != "" {
			err = global.DB.Model(dbTag).Preload("Tags").Order(sortField[info.Sort] + " desc").Association("Articles").Find(&list)
		} else {
			err = global.DB.Model(dbTag).Preload("Tags").Order("created_at desc").Association("Articles").Find(&list)
		}

	}

	if info.Name == "articles" {
		if info.Sort != "" {
			err = db.Where("title like ?", strings.Join([]string{"%", text, "%"}, "")).Preload("Tags").Order(sortField[info.Sort] + " desc").Find(&list).Error
		} else {
			err = db.Where("title like ?", strings.Join([]string{"%", text, "%"}, "")).Preload("Tags").Order("created_at desc").Find(&list).Error
		}
	}

	return
}
