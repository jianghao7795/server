package app

import (
	"errors"

	global "server/model"

	"gorm.io/gorm"
)

type PraiseService struct{}

var PraiseServer = new(PraiseService)

// LikeComment 点赞评论
func (*PraiseService) LikeComment(commentId uint, userId uint) (*global.Praise, error) {
	praise := &global.Praise{CommentId: commentId, UserId: userId}
	err := global.DB.Create(praise).Error
	if err == nil {
		return praise, nil
	}

	if !errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, err
	}

	var existing global.Praise
	if err := global.DB.Unscoped().Where("user_id = ? AND comment_id = ?", userId, commentId).First(&existing).Error; err != nil {
		return nil, err
	}

	if existing.DeletedAt.Valid {
		global.DB.Unscoped().Model(&existing).Update("deleted_at", nil)
		existing.DeletedAt = gorm.DeletedAt{}
	}
	return &existing, nil
}

// UnlikeComment 取消点赞评论
func (*PraiseService) UnlikeComment(commentId uint, userId uint) error {
	result := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&global.Praise{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到点赞记录")
	}
	return nil
}

// CheckCommentLiked 检查用户是否已点赞
func (*PraiseService) CheckCommentLiked(commentId uint, userId uint) (bool, *global.Praise, error) {
	var praise global.Praise
	err := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).First(&praise).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &praise, nil
}
