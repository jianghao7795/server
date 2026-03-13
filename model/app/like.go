package app

import (
	global "server/model"
)

// Like 点赞结构体
type Like struct {
	global.MODEL
	UserId uint `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户ID;not null;index"`
	PostId uint `json:"post_id" form:"post_id" gorm:"column:post_id;comment:帖子ID;not null;index"`
	User   User `json:"user" form:"user" gorm:"foreignKey:UserId"`
	Post   Post `json:"post" form:"post" gorm:"foreignKey:PostId"`
}

// TableName Like 表名
func (Like) TableName() string {
	return "likes"
}
