package system

import (
	global "server/model"
)

type SysUserProblem struct {
	global.MODEL
	SysUserId int    `query:"sys_user_id" json:"sys_user_id" form:"sys_user_id" gorm:"column:sys_user_id;comment:用户的ID"`
	Problem   string `json:"problem" form:"problem" gorm:"column:problem;comment:问题"`
	Answer    string `json:"answer" form:"answer" gorm:"column:answer;comment:答案"`
}

// TableName Comment 表名
func (SysUserProblem) TableName() string {
	return "sys_user_problems"
}
