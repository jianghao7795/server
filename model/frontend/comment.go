package frontend

import global "server/model"

// Comment 结构体
// 如果含有time.Time 请自行import time包
type Comment struct {
	global.MODEL
	ArticleId int       `query:"article_id" json:"article_id" form:"article_id" gorm:"column:article_id;comment:文章id;size:10;"`
	Article   Article   `json:"article" gorm:"foreignKey:ArticleId"`
	ParentId  int       `query:"parent_id" json:"parent_id" form:"parent_id" gorm:"column:parent_id;comment:上级;size:10;"`
	Content   string    `json:"content" form:"content" gorm:"column:content;comment:内容;"`
	UserId    int       `query:"user_id" json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	User      User      `json:"user" form:"user" gorm:"foreignKey:UserId"`
	ToUserId  int       `query:"to_user_id" json:"to_user_id" form:"to_user_id" gorm:"column:to_user_id;comment:回复评论的用户id;"`
	ToUser    User      `query:"to_user" json:"to_user" form:"to_user" gorm:"foreignKey:ToUserId"`
	Children  []Comment `json:"children" form:"children" gorm:"foreignKey:ParentId;"`
}

// TableName Comment 表名
func (Comment) TableName() string {
	return "comments"
}
