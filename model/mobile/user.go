package mobile

import (
	global "server/model"
)

// MobileUser 结构体
// 如果含有time.Time 请自行import time包
type MobileUser struct {
	global.MODEL
	Username string `json:"username" form:"username" gorm:"column:username;comment:用户名;size:50;"`
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;size:50;"`
	Realname string `json:"realname" form:"realname" gorm:"column:realname;comment:真实姓名;size:50;"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"column:avatar;comment:头像;size:255;"`
	Sign     string `json:"sign" form:"sign" gorm:"column:sign;comment:简介;size:255;"`
	Cover    string `json:"cover" form:"cover" gorm:"column:cover;comment:主页封面;size:255;"`
	Content  string `json:"content" form:"content" gorm:"column:content;comment:用户信息;"`
	// Password string `json:"password" form:"password" gorm:"column:password;comment:密码;size:100;"`
	Industry uint8  `json:"industry" form:"industry" gorm:"column:industry;comment:行业;"`
	Gender   uint8  `json:"gender" form:"gender" gorm:"column:gender;comment:性别;"`
	Phone    string `json:"phone" form:"phone" gorm:"column:phone;comment:电话;"`
}

// TableName MobileUser 表名
func (MobileUser) TableName() string {
	return "mobile_users"
}
