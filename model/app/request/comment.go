package request

import (
	"server/model/common/request"
)

type CommentSearch struct {
	ArticleId int    `json:"articleId" form:"articleId"`
	Content   string `json:"content" form:"content"`
	request.PageInfo
}

// 点赞参数
type CommentParise struct {
	ID     int `json:"id" form:"id"`
	Parise int `json:"parise" form:"parise"`
}
