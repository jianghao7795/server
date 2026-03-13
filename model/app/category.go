package app

import (
	global "server/model"
)

// Category 分类结构体
type Category struct {
	global.MODEL
	Name string `json:"name" form:"name" gorm:"column:name;comment:分类名称;size:100;not null"`
	Sort int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0"`
}

// TableName Category 表名
func (Category) TableName() string {
	return "categories"
}
