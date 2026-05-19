package frontend

import (
	"errors"

	global "server/model"
	"server/model/frontend"

	"gorm.io/gorm"
)

type Comment struct{}

func (commentService *Comment) GetCommentByArticleId(articleId int) (list []frontend.Comment, err error) {
	db := global.DB.Model(&frontend.Comment{})
	var commentList []frontend.Comment
	err = db.Where("article_id = ?", articleId).Where("parent_id = ?", 0).Preload("User").Preload("Praises").Order("id desc").Find(&commentList).Error
	if len(commentList) > 0 {
		for comment := range commentList {
			err = commentService.findChildrenComment(&commentList[comment])
			if err != nil {
				break
			}
		}
	}

	return commentList, err
}

func (s *Comment) findChildrenComment(comment *frontend.Comment) (err error) {
	err = global.DB.Where("parent_id = ?", comment.ID).Preload("User").Preload("ToUser").Preload("Praises").Order("id asc").Find(&comment.Children).Error
	for i := range comment.Children {
		if err = s.findChildrenComment(&comment.Children[i]); err != nil {
			return err
		}
	}
	return nil
}

func (*Comment) CreatedComment(info *frontend.Comment) (err error) {
	err = global.DB.Create(info).Error
	return
}

// LikeComment 点赞评论
func (*Comment) LikeComment(commentId uint, userId int64) (*frontend.Praise, error) {
	var existing frontend.Praise
	err := global.DB.Unscoped().Where("user_id = ? AND comment_id = ?", userId, commentId).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		praise := &frontend.Praise{CommentId: int64(commentId), UserId: userId}
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
func (*Comment) UnlikeComment(commentId uint, userId int64) error {
	result := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&frontend.Praise{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到点赞记录")
	}
	return nil
}

// CheckCommentLiked 检查用户是否已点赞
func (*Comment) CheckCommentLiked(commentId uint, userId int64) (bool, *frontend.Praise, error) {
	var praise frontend.Praise
	err := global.DB.Where("user_id = ? AND comment_id = ?", userId, commentId).First(&praise).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &praise, nil
}
