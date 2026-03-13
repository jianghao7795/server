package request

import (
	"server/model/common/request"
	"server/model/mobile"
)

type MobileUserSearch struct {
	mobile.MobileUser
	request.PageInfo
}

type MobileUpdate struct {
	Field string `json:"field" form:"field"`
	Value any    `json:"value" form:"value"`
}

// update password struct
type MobileUpdatePassword struct {
	ID          uint   `json:"id" form:"id"`
	Password    string `json:"password" form:"password"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}
