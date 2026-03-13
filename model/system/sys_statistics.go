package system

import global "server/model"

type SysStatistics struct {
	global.MODEL
	Tablename string `json:"table_name" form:"table_name" grom:"comment:表名"`
	Status    uint   `json:"status" form:"status" grom:"comment:表的状态"`
	Total     uint   `json:"total" form:"total" grom:"comment:总数"`
}

func (SysStatistics) TableName() string {
	return "sys_statistics"
}
