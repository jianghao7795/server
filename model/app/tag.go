// 自动生成模板Tag
package app

import (
	global "server/model"
)

// Tag 结构体
// 如果含有time.Time 请自行import time包
type Tag struct {
	global.MODEL
	Name     string    `json:"name" form:"name" gorm:"column:name;comment:标签名称;size:191;"`
	Status   int       `json:"status" form:"status" gorm:"column:status;comment:状态;"`
	Articles []Article `json:"articles" form:"articles" gorm:"many2many:article_tag"`
}

// TableName Tag 表名
func (Tag) TableName() string {
	return "tags"
}
