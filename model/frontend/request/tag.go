package request

import (
	"server/model/app"
	"server/model/common/request"
)

type TagSearch struct {
	app.Tag
	request.PageInfo
}
