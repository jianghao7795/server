package app

import (
	"errors"
	global "server/model"
	"server/model/app"

	"gorm.io/gorm"
)

type LikeService struct{}

// LikePost 点赞帖子
func (*LikeService) LikePost(postId, userId uint) error {
	// 检查是否已经点赞
	var existingLike app.Like
	err := global.DB.Where("post_id = ? AND user_id = ?", postId, userId).First(&existingLike).Error
	if err == nil {
		return errors.New("您已经点赞过这篇帖子了")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 检查帖子是否存在
	var post app.Post
	if err := global.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	// 开始事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建点赞记录
		like := app.Like{
			PostId: postId,
			UserId: userId,
		}
		if err := tx.Create(&like).Error; err != nil {
			return err
		}

		// 更新帖子点赞数
		if err := tx.Model(&app.Post{}).Where("id = ?", postId).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// UnlikePost 取消点赞帖子
func (*LikeService) UnlikePost(postId, userId uint) error {
	// 检查是否已经点赞
	var existingLike app.Like
	err := global.DB.Where("post_id = ? AND user_id = ?", postId, userId).First(&existingLike).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您还没有点赞过这篇帖子")
		}
		return err
	}

	// 开始事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Delete(&existingLike).Error; err != nil {
			return err
		}

		// 更新帖子点赞数
		if err := tx.Model(&app.Post{}).Where("id = ?", postId).Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPostLikes 获取帖子点赞列表
func (*LikeService) GetPostLikes(postId uint, page, pageSize int) ([]app.Like, int64, error) {
	var likes []app.Like
	var total int64

	// 计算总数
	if err := global.DB.Model(&app.Like{}).Where("post_id = ?", postId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := global.DB.Where("post_id = ?", postId).
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&likes).Error

	return likes, total, err
}

// CheckUserLiked 检查用户是否点赞了某篇帖子
func (*LikeService) CheckUserLiked(postId, userId uint) (bool, error) {
	var count int64
	err := global.DB.Model(&app.Like{}).Where("post_id = ? AND user_id = ?", postId, userId).Count(&count).Error
	return count > 0, err
}

// GetUserLikedPosts 获取用户点赞的帖子列表
func (*LikeService) GetUserLikedPosts(userId uint, page, pageSize int) ([]app.Like, int64, error) {
	var likes []app.Like
	var total int64

	// 计算总数
	if err := global.DB.Model(&app.Like{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := global.DB.Where("user_id = ?", userId).
		Preload("Post").
		Preload("Post.User").
		Preload("Post.Categories").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&likes).Error

	return likes, total, err
}

// GetPostLikeCount 获取帖子点赞数
func (*LikeService) GetPostLikeCount(postId uint) (int64, error) {
	var count int64
	err := global.DB.Model(&app.Like{}).Where("post_id = ?", postId).Count(&count).Error
	return count, err
}
