package frontend

import (
	"encoding/json"
	"errors"
	global "server/model"
	"server/model/frontend"
	"strconv"
	"time"

	// json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"

	"github.com/redis/go-redis/v9"
)

type Tag struct{}

func (s *Tag) GetTagList(c fiber.Ctx) (list []frontend.Tag, err error) {
	var tagListStr string
	tagListStr, err = global.REDIS.Get(c.Context(), "tag-list").Result()
	cacheTime := global.CONFIG.Cache.Time
	if errors.Is(err, redis.Nil) {
		db := global.DB.Model(&frontend.Tag{})
		err = db.Order("id desc").Find(&list).Error
		if err != nil {
			return list, err
		}
		strList, _ := json.Marshal(list)
		err := global.REDIS.Set(c.Context(), "tag-list", strList, time.Duration(cacheTime)*time.Second).Err()
		if err != nil {
			return list, errors.New("redis 存储失败， 请查看redis服务 配置")
		}
		return list, nil
	} else if err != nil {
		return list, err
	} else {
		if tagListStr != "" {
			err = json.Unmarshal([]byte(tagListStr), &list)
			return list, err
		}
	}
	return list, err
}

func (s *Tag) GetTagArticle(tagId int, c fiber.Ctx) (tagArticle frontend.Tag, err error) {
	tagArticleStr := ""
	tagArticleStr, err = global.REDIS.Get(c.Context(), "tag-"+strconv.Itoa(tagId)).Result()
	if errors.Is(err, redis.Nil) {
		db := global.DB.Model(&frontend.Tag{})
		err = db.Where("id = ?", tagId).Order("id desc").Preload("Articles").First(&tagArticle).Error
		if err != nil {
			return tagArticle, err
		}
		err := global.REDIS.Set(c.Context(), "tag-list", tagArticle, 1*time.Hour).Err()
		if err != nil {
			return tagArticle, errors.New("redis 存储失败， 请查看redis服务 配置")
		}
		return tagArticle, nil
		// return tagArticle, err
	} else if err != nil {
		return tagArticle, err
	} else {
		err = json.Unmarshal([]byte(tagArticleStr), &tagArticle)
	}

	return tagArticle, err
}
