package request

import (
	"server/model/common/request"
	"server/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
