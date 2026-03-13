// 自动生成模板User
package app

import (
	global "server/model"
)

// User 结构体
// 如果含有time.Time 请自行import time包
type User struct {
	global.MODEL
	Name         string `json:"name" form:"name" gorm:"column:name;comment:用户名;size:50;"`
	HeadImg      string `json:"headImg" form:"headImg" gorm:"column:head_img;comment:背景图;"`
	Introduction string `json:"introduction" form:"introduction" gorm:"column:introduction;comment:简介;size:255;"`
	Content      string `json:"content" form:"content" gorm:"column:content;comment:用户信息;size:255;"`
	Password     string `json:"password" form:"password" gorm:"column:password;comment:密码;size:255;"`
	Header       string `json:"header" form:"header" gorm:"column:header;comment:头像;size:100"`
}

type UserFrontend struct {
	global.MODEL
	Name         string `json:"name" form:"name" gorm:"column:name;comment:用户名;size:50;"`
	HeadImg      string `json:"headImg" form:"headImg" gorm:"column:head_img;comment:背景图;"`
	Introduction string `json:"introduction" form:"introduction" gorm:"column:introduction;comment:简介;size:255;"`
	Content      string `json:"content" form:"content" gorm:"column:content;comment:用户信息;size:255;"`
	Header       string `json:"header" form:"header" gorm:"column:header;comment:头像;size:100"`
}

// TableName User 表名
func (User) TableName() string {
	return "users"
}

func (UserFrontend) TableName() string {
	return "users"
}
