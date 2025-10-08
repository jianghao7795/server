package app

import (
	"errors"
	"fmt"
	global "server-fiber/model"
	"server-fiber/model/app"
	appReq "server-fiber/model/app/request"
	"server-fiber/model/common/request"
	"server-fiber/utils"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	*utils.CRUDBase[app.Article]
}

// NewArticleService creates a new article service
func NewArticleService() *ArticleService {
	return &ArticleService{
		CRUDBase: utils.NewCRUDBase[app.Article](global.DB),
	}
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(article *app.Article) error {
	return s.Create(article)
}

// DeleteArticle deletes an article and its related comments
func (s *ArticleService) DeleteArticle(id uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// Delete related comments first
		if err := tx.Delete(&app.Comment{}, "article_id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete related comments: %w", err)
		}

		// Delete the article
		if err := tx.Delete(&app.Article{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete article: %w", err)
		}

		return nil
	})
}

// DeleteArticleByIds deletes multiple articles by IDs
func (s *ArticleService) DeleteArticleByIds(ids request.IdsReq) error {
	return s.DeleteByIDs(ids)
}

// UpdateArticle updates an existing article
func (s *ArticleService) UpdateArticle(article *app.Article) error {
	return s.Update(article)
}

// GetArticle retrieves an article by ID with related data
func (s *ArticleService) GetArticle(id uint) (app.Article, error) {
	var article app.Article
	err := global.DB.Preload("Tags").Preload("User").Where("id = ?", id).First(&article).Error
	return article, err
}

// getList

func (*ArticleService) GetArticleInfoList(info *appReq.ArticleSearch) (list []app.Article, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&app.Article{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Title != "" {
		db = db.Where("title like ?", "%"+info.Title+"%")
	}
	if info.IsImportant != 0 {
		db = db.Where("is_important = ?", info.IsImportant)
	}
	//
	err = db.Limit(limit).Offset(offset).Order("id desc").Preload("User").Preload("Tags").Find(&list).Error
	return list, total, err
}

// 批量更新
func (*ArticleService) PutArticleByIds(ids *request.IdsReq) (err error) {
	result := global.DB.Model(&app.Article{}).Where("id in ?", ids.Ids).Update("is_important", 2)
	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected == 0 {
		err = errors.New("没有更新任何数据")
		return
	}
	return
}

func (*ArticleService) GetArticleReading(userId uint) (count int64, err error) {
	t := time.Now()
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 99, t.Location())
	err = global.DB.Model(&app.Ip{}).Where("user_id = ?", userId).Where("created_at > ? and created_at < ?", startTime, endTime).Count(&count).Error
	return
}
