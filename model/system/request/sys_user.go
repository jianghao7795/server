package request

import (
	"server/model/common/request"
	model "server/model/system"
)

// User register structure
type Register struct {
	Username     string   `json:"userName"`
	Password     string   `json:"passWord"`
	NickName     string   `json:"nickName" gorm:"default:'QMPlusUser'"`
	HeaderImg    string   `json:"headerImg" gorm:"default:'public/logo.png'"`
	AuthorityId  string   `json:"authorityId" gorm:"default:888"`
	AuthorityIds []string `json:"authorityIds"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
}

// User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// get token login struct
type LoginToken struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// Modify password structure
type ChangePasswordStruct struct {
	Username    string `json:"username"`    // 用户名
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId string `json:"authorityId"` // 角色ID
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []string `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID           uint                 `gorm:"primarykey"`                                            // 主键ID
	NickName     string               `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`             // 用户昵称
	Phone        string               `json:"phone"  gorm:"comment:用户手机号"`                           // 用户角色ID
	AuthorityIds []string             `json:"authorityIds" gorm:"-"`                                 // 角色ID
	Email        string               `json:"email"  gorm:"comment:用户邮箱"`                            // 用户邮箱
	HeaderImg    string               `json:"headerImg" gorm:"default:public/logo.png;comment:用户头像"` // 用户头像
	Authorities  []model.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}

type SearchInfo struct {
	request.PageInfo
	Username string `json:"username" form:"username"`
}
