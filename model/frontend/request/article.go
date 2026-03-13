package request

import (
	"server/model/app"
	"server/model/common/request"
)

type ArticleSearch struct {
	app.Article
	request.PageInfo
	Name  string `json:"name" form:"name"`
	Value string `json:"value" form:"value"`
	Sort  string `json:"sort" form:"sort"`
}
