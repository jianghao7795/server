package app

import (
	"errors"
	global "server/model"
	"server/model/app"
	commentReq "server/model/app/request"
	"server/model/common/request"
	"strings"
)

type CommentService struct{}

func (commentService *CommentService) CreateComment(comment *app.Comment) (err error) {
	err = global.DB.Create(comment).Error
	return err
}

func (commentService *CommentService) DeleteComment(id uint) (err error) {
	err = global.DB.Delete(&app.Comment{}, id).Error
	return err
}

func (commentService *CommentService) DeleteCommentByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]app.Comment{}, "id in ?", ids.Ids).Error
	return err
}

func (commentService *CommentService) UpdateComment(comment *app.Comment) (err error) {
	var commentReplica app.Comment
	if err = global.DB.Where("id = ?", comment.ID).First(&commentReplica).Error; err != nil {
		return errors.New("未找到该comment")
	}
	if commentReplica.ID == 0 {
		return errors.New("未找到该comment")
	}
	result := global.DB.Save(comment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("没有更新任何数据")
	}
	return nil
}

func (commentService *CommentService) GetComment(id int) (comment app.Comment, err error) {
	err = global.DB.Preload("Article").Preload("Praises").Where("id = ?", id).First(&comment).Error
	return
}

func (commentService *CommentService) GetCommentList(postId uint, page, pageSize int) ([]app.Comment, int64, error) {
	var comments []app.Comment
	var total int64

	if err := global.DB.Model(&app.Comment{}).Where("post_id = ?", postId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := global.DB.Where("post_id = ?", postId).
		Preload("User").
		Preload("ToUser").
		Preload("Praises").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error

	return comments, total, err
}

func (commentService *CommentService) GetCommentInfoList(info *commentReq.CommentSearch) (list []app.Comment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&app.Comment{}).Preload("Article").Preload("User").Preload("ToUser").Preload("Praises")
	if info.ArticleId != 0 {
		db = db.Where("article_id = ?", info.ArticleId)
	}
	if info.Content != "" {
		db = db.Where("content like ?", strings.Join([]string{"%", info.Content, "%"}, ""))
	}
	var comments []app.Comment
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&comments).Error
	return comments, total, err
}

func (commentService *CommentService) GetCommentTreeList(info *commentReq.CommentSearch) (list []app.Comment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&app.Comment{})
	if info.ArticleId != 0 {
		db = db.Where("article_id = ?", info.ArticleId)
	}
	if info.Content != "" {
		db = db.Where("content like ?", strings.Join([]string{"%", info.Content, "%"}, ""))
	}
	err = db.Where("parent_id = ?", 0).Count(&total).Error
	if err != nil {
		return
	}

	var commentList []app.Comment
	err = db.Limit(limit).Offset(offset).Where("parent_id = ?", 0).Preload("Article").Preload("User").Preload("Praises").Order("id desc").Find(&commentList).Error
	if len(commentList) > 0 {
		for comment := range commentList {
			err = commentService.findChildrenComment(&commentList[comment])
		}
	}
	return commentList, total, err
}

func (commentService *CommentService) findChildrenComment(comment *app.Comment) (err error) {
	err = global.DB.Where("parent_id = ?", comment.ID).Preload("User").Preload("ToUser").Preload("Praises").Order("id asc").Find(&comment.Children).Error
	for i := range comment.Children {
		if err = commentService.findChildrenComment(&comment.Children[i]); err != nil {
			return err
		}
	}
	return nil
}
