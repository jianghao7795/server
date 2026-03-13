package request

import (
	"server/model/app"
	"server/model/common/request"
)

type BaseMessageBody struct {
	app.BaseMessage
	request.PageInfo
}
