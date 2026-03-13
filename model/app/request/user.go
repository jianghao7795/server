package request

import (
	"server/model/app"
	"server/model/common/request"
)

type UserSearch struct {
	app.User
	request.PageInfo
}
