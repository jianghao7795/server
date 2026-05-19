package app

import (
	"errors"
	global "server/model"
	"server/model/app"
	commentReq "server/model/app/request"
	"server/model/common/request"
	"strings"

	"gorm.io/gorm"
)

type CommentService struct{}

// CreateComment 创建Comment记录
func (commentService *CommentService) CreateComment(comment *app.Comment) (err error) {
	err = global.DB.Create(comment).Error
	return err
}

// DeleteComment 删除Comment记录
func (commentService *CommentService) DeleteComment(id uint) (err error) {
	err = global.DB.Delete(&app.Comment{}, id).Error
	return err
}

// DeleteCommentByIds 批量删除Comment记录
func (commentService *CommentService) DeleteCommentByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]app.Comment{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateComment 更新Comment记录
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
func (commentService *CommentService) GetComment(id int) (comment app.Comment, err error) {
	err = global.DB.Preload("Article").Preload("Praises").Where("id = ?", id).First(&comment).Error
	return
}

// GetCommentList 根据帖子ID获取评论列表
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

// GetCommentInfoList 分页获取Comment记录
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

// GetCommentTreeList 分页获取TreeList
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

// PutLikeItOrDislike 评论点赞/取消点赞
func (*CommentService) PutLikeItOrDislike(info *app.Praise) (err error) {
	if info.ID == 0 {
		// 查找已有记录（包括软删除的）
		var existing app.Praise
		err = global.DB.Unscoped().Where("user_id = ? AND comment_id = ?", info.UserId, info.CommentId).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 首次点赞 → 新建
				return global.DB.Create(info).Error
			}
			return err
		}

		// 找到已有记录
		info.ID = existing.ID
		if existing.DeletedAt.Valid {
			// 曾经取消过点赞 → 恢复
			return global.DB.Unscoped().Model(&existing).Update("deleted_at", nil).Error
		}
		// 已点赞 → 返回已有记录
		info.CreatedAt = existing.CreatedAt
		info.UpdatedAt = existing.UpdatedAt
		return nil
	}

	// 取消点赞
	return global.DB.Where("id = ?", info.ID).Delete(&app.Praise{}).Error
}

// LikeComment 点赞评论
func (*CommentService) LikeComment(commentId uint, userId int64) (*app.Praise, error) {
	var existing app.Praise
	err := global.DB.Unscoped().Where("user_id = ? AND comment_id = ?", userId, commentId).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		praise := &app.Praise{CommentId: int64(commentId), UserId: userId}
		return praise, global.DB.Create(praise).Error
	}
	if err != nil {
		return nil, err
	}

	if existing.DeletedAt.Valid {
		global.DB.Unscoped().Model(&existing).Update("deleted_at", nil)
		existing.DeletedAt = gorm.DeletedAt{}
	}
	return &existing, nil
}

// UnlikeComment 取消点赞评论
func (*CommentService) UnlikeComment(commentId uint, userId int64) error {
	result := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&app.Praise{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到点赞记录")
	}
	return nil
}

// CheckCommentLiked 检查用户是否已点赞
func (*CommentService) CheckCommentLiked(commentId uint, userId int64) (bool, *app.Praise, error) {
	var praise app.Praise
	err := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).First(&praise).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &praise, nil
}
