package request

import (
	"server/model/app"
	"server/model/common/request"
)

type ArticleSearch struct {
	app.Article
	request.PageInfo
}
