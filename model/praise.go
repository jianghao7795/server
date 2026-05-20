package model

type Praise struct {
	MODEL
	CommentId uint `json:"comment_id" form:"comment_id" query:"comment_id" gorm:"column:comment_id;comment:评论id;size:10;uniqueIndex:uk_praise;not null"`
	UserId    uint `json:"user_id" form:"user_id" query:"user_id" gorm:"column:user_id;comment:用户id;size:10;uniqueIndex:uk_praise;not null"`
}

func (Praise) TableName() string {
	return "praise"
}
