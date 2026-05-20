package app

import (
	global "server/model"
	"server/model/system"
)

type Comment struct {
	global.MODEL
	PostId    uint           `json:"post_id" form:"post_id" gorm:"column:post_id;comment:帖子ID;not null;index"`
	Post      Post           `json:"post" gorm:"foreignKey:PostId"`
	ParentId  uint           `json:"parent_id" form:"parent_id" query:"parent_id" gorm:"column:parent_id;comment:父评论ID;default:0;index"`
	Content   string         `json:"content" form:"content" gorm:"column:content;comment:评论内容;type:text;"`
	UserId    uint           `json:"user_id" form:"user_id" query:"user_id" gorm:"column:user_id;comment:用户ID;not null;index"`
	User      system.SysUser `json:"user" form:"user" gorm:"foreignKey:UserId"`
	Children  []Comment      `json:"children" form:"children" gorm:"foreignKey:ParentId"`
	ToUserId  uint           `json:"to_user_id" form:"to_user_id" query:"to_user_id" gorm:"column:to_user_id;comment:回复用户ID;default:0;index"`
	ToUser    system.SysUser `json:"to_user" form:"to_user" gorm:"foreignKey:ToUserId"`
	ArticleId uint           `json:"article_id" form:"article_id" query:"article_id" gorm:"column:article_id;comment:文章ID;not null;index"`
	Article   Article        `json:"article" gorm:"foreignKey:ArticleId"`
	Praises   []global.Praise `json:"praises" gorm:"foreignKey:CommentId"`
}

func (Comment) TableName() string {
	return "comments"
}
