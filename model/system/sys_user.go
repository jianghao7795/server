package system

import (
	global "server/model"

	uuid "github.com/google/uuid"
)

type SysUser struct {
	global.MODEL
	UUID         uuid.UUID      `json:"uuid" gorm:"comment:用户UUID"`                            // 用户UUID
	Username     string         `json:"userName" gorm:"comment:用户登录名"`                         // 用户登录名
	Password     string         `json:"-"  gorm:"comment:用户登录密码"`                              // 用户登录密码
	NickName     string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`             // 用户昵称
	SideMode     string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`           // 用户侧边主题
	HeaderImg    string         `json:"headerImg" gorm:"default:public/logo.png;comment:用户头像"` // 用户头像
	BaseColor    string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`            // 基础颜色
	ActiveColor  string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`       // 活跃颜色
	AuthorityId  string         `json:"authorityId" gorm:"default:888;comment:用户角色ID"`         // 用户角色ID
	Authority    SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities  []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone        string         `json:"phone"  gorm:"comment:用户手机号"` // 用户手机号
	Email        string         `json:"email"  gorm:"comment:用户邮箱"`  // 用户邮箱
	HeadImg      string         `query:"head_img" json:"head_img" gorm:"comment:背景图"`
	Introduction string         `json:"introduction" gorm:"comment:简介"`
	Content      string         `json:"content" gorm:"comment:介绍"`
}

func (SysUser) TableName() string {
	return "sys_users"
}
