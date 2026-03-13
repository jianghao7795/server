package app

import (
	global "server/model"
	"server/model/system"
)

// Post 帖子结构体
type Post struct {
	global.MODEL
	Title           string       `json:"title" form:"title" gorm:"column:title;comment:帖子标题;size:191;"`
	Text            string       `json:"text" form:"text" gorm:"column:text;comment:帖子内容;"`
	PublishAt       global.MODEL `json:"publish_at" form:"publish_at" gorm:"column:publish_at;comment:发布时间;"`
	UserId          uint         `json:"user_id" form:"user_id" query:"user_id" gorm:"column:user_id;comment:帖子作者ID;not null;index"`
	State           int          `json:"state" form:"state" gorm:"column:state;comment:帖子状态;default:1"`
	IsImportant     int          `json:"is_important" form:"is_important" query:"is_important" gorm:"column:is_important;comment:首页是否显示;default:0"`
	ReadingQuantity int          `json:"reading_quantity" form:"reading_quantity" query:"reading_quantity" gorm:"column:reading_quantity;comment:阅读量;default:0"`
	LikeCount       int          `json:"like_count" form:"like_count" query:"like_count" gorm:"column:like_count;comment:点赞数;default:0"`

	// 关联关系
	User       system.SysUser `json:"user" form:"user" gorm:"foreignKey:UserId"`
	Comments   []Comment      `json:"comments" form:"comments" gorm:"foreignKey:PostId"`
	Likes      []Like         `json:"likes" form:"likes" gorm:"foreignKey:PostId"`
	Categories []Category     `json:"categories" form:"categories" gorm:"many2many:category_posts"`
}

// TableName Post 表名
func (Post) TableName() string {
	return "posts"
}
