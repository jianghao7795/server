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

// CreateComment 创建Comment记录
// Author wuhao
func (commentService *CommentService) CreateComment(comment *app.Comment) (err error) {
	err = global.DB.Create(comment).Error
	return err
}

// DeleteComment 删除Comment记录
// Author wuhao
func (commentService *CommentService) DeleteComment(id uint) (err error) {
	err = global.DB.Delete(&app.Comment{}, id).Error
	return err
}

// DeleteCommentByIds 批量删除Comment记录
// Author wuhao
func (commentService *CommentService) DeleteCommentByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]app.Comment{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateComment 更新Comment记录
// Author wuhao
func (commentService *CommentService) UpdateComment(comment *app.Comment) (err error) {
	var commentReplica app.Comment
	db := global.DB.Where("id = ?", comment.ID).First(&commentReplica)
	if commentReplica.ID == 0 {
		return errors.New("未找到该comment")
	}
	result := db.Save(comment)
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

// GetComment 根据id获取Comment记录
// Author wuhao
func (commentService *CommentService) GetComment(id int) (comment app.Comment, err error) {
	err = global.DB.Preload("Article").Where("id = ?", id).First(&comment).Error
	return
}

// GetCommentList 根据帖子ID获取评论列表
func (commentService *CommentService) GetCommentList(postId uint, page, pageSize int) ([]app.Comment, int64, error) {
	var comments []app.Comment
	var total int64

	// 计算总数
	if err := global.DB.Model(&app.Comment{}).Where("post_id = ?", postId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := global.DB.Where("post_id = ?", postId).
		Preload("User").
		Preload("ToUser").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error

	return comments, total, err
}

// GetCommentInfoList 分页获取Comment记录
// Author wuhao
func (commentService *CommentService) GetCommentInfoList(info *commentReq.CommentSearch) (list []app.Comment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&app.Comment{}).Preload("Article").Preload("User").Preload("ToUser").Preload("Praise")
	if info.ArticleId != 0 {
		db = db.Where("article_id = ?", info.ArticleId)
	}
	if info.Content != "" {
		db = db.Where("content like ?", strings.Join([]string{"%", info.Content, "%"}, ""))
		// db = db.Where("MATCH(content) AGAINST('+" + info.Content + "')")
	}
	var comments []app.Comment
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&comments).Error
	return comments, total, err
}

// GetCommentTreeList 分页获取Treelist

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

	err = db.Limit(limit).Offset(offset).Where("parent_id = ?", 0).Preload("Article").Preload("User").Order("id desc").Find(&commentList).Error
	// err = db.Limit(limit).Offset(offset).Where("parent_id = ?", 0).Find(&commentList).Error
	if len(commentList) > 0 {
		for comment := range commentList {
			err = commentService.findChildrenComment(&commentList[comment])
		}
	}
	return commentList, total, err
}

func (commentService *CommentService) findChildrenComment(comment *app.Comment) (err error) {
	err = global.DB.Where("parent_id = ?", comment.ID).Preload("User").Preload("ToUser").Order("user_id desc").Find(&comment.Children).Error
	return err
}

// LikeIt 点赞一条记录
func (*CommentService) PutLikeItOrDislike(info *app.Praise) (err error) {
	db := global.DB.Model(&app.Praise{})

	if info.ID == 0 {
		var praise app.Praise
		err = db.Raw("Select id, comment_id, user_id, created_at, updated_at from praise where user_id = ? and comment_id = ? limit 1", info.UserId, info.CommentId).Scan(&praise).Error
		if err != nil {
			return err
		}
		if praise.ID == 0 {
			err = db.Create(info).Error
		} else {
			info.ID = praise.ID
			info.CreatedAt = praise.CreatedAt
			info.UpdatedAt = praise.UpdatedAt
			err = db.Exec("UPDATE `praise` SET `deleted_at`=NULL,`comment_id`=?,`user_id`=? where id = ? ORDER BY `praise`.`id` LIMIT 1", info.CommentId, info.UserId, info.ID).Error

		}
	} else {
		err = db.Where("id = ?", info.ID).Delete(info).Error
	}

	// if err != nil {
	// 	return err
	// }
	return err
}
