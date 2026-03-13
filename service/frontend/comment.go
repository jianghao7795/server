package frontend

import (
	global "server/model"
	"server/model/frontend"
)

type Comment struct{}

func (commentService *Comment) GetCommentByArticleId(articleId int) (list []frontend.Comment, err error) {
	db := global.DB.Model(&frontend.Comment{})
	var commentList []frontend.Comment
	err = db.Where("article_id = ?", articleId).Where("parent_id = ?", 0).Preload("User").Order("id desc").Find(&commentList).Error
	// err = db.Limit(limit).Offset(offset).Where("parent_id = ?", 0).Find(&commentList).Error
	if len(commentList) > 0 {
		for comment := range commentList {
			err = commentService.findChildrenComment(&commentList[comment])
		}
	}

	return commentList, err
}

func (*Comment) findChildrenComment(comment *frontend.Comment) (err error) {
	err = global.DB.Where("parent_id = ?", comment.ID).Preload("User").Preload("ToUser").Order("user_id desc").Find(&comment.Children).Error
	return err
}

func (*Comment) CreatedComment(info *frontend.Comment) (err error) {
	err = global.DB.Create(info).Error
	return
}
