package app

import (
	global "server/model"
)

// CategoryPost 分类帖子关联结构体
type CategoryPost struct {
	global.MODEL
	CategoryId uint     `json:"category_id" form:"category_id" gorm:"column:category_id;comment:分类ID;not null;index"`
	PostId     uint     `json:"post_id" form:"post_id" gorm:"column:post_id;comment:帖子ID;not null;index"`
	Category   Category `json:"category" form:"category" gorm:"foreignKey:CategoryId"`
	Post       Post     `json:"post" form:"post" gorm:"foreignKey:PostId"`
}

// TableName CategoryPost 表名
func (CategoryPost) TableName() string {
	return "category_posts"
}
