package app

import global "server/model"

type Ip struct {
	global.MODEL
	Ip        string `json:"ip" form:"ip" gorm:"column:ip;comment:ip;size:50;"`
	ArticleID uint   `json:"article_id" form:"article_id" query:"article_id" gorm:"column:article_id;comment:文章id"`
	UserID    uint   `json:"user_id" form:"user_id" query:"user_id" gorm:"column:user_id;comment:用户id"`
}

func (Ip) TableName() string {
	return "ips"
}
