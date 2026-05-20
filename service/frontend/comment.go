package frontend

import (
	appService "server/service/app"

	global "server/model"
	"server/model/frontend"
)

type Comment struct{}

var praiseService = appService.PraiseServer

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
