package request

import (
	"server/model/common/request"
)

// api分页条件查询及排序结构体
type SearchCustomerParams struct {
	CustomerName      string `json:"customerName" form:"customerName" gorm:"comment:客户名"`
	CustomerPhoneData string `json:"customerPhoneData" form:"customerPhoneData" gorm:"comment:客户手机号"`
	request.PageInfo
}
