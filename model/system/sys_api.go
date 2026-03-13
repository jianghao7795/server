package system

import (
	global "server/model"
)

type SysApi struct {
	global.MODEL
	Path        string `json:"path" form:"path" gorm:"comment:api路径"`                            // api路径
	Description string `json:"description" form:"description" gorm:"comment:api中文描述"`            // api中文描述
	ApiGroup    string `query:"api_group" json:"api_group" form:"api_group" gorm:"comment:api组"` // api组
	Method      string `json:"method" form:"method" gorm:"default:POST;comment:方法"`              // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}
