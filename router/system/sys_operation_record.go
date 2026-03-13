package system

import (
	v1 "server/api/v1/system"

	"github.com/gofiber/fiber/v3"
)

type OperationRecordRouter struct{}

func (s *OperationRecordRouter) InitSysOperationRecordRouter(Router fiber.Router) {
	operationRecordRouter := Router.Group("sysOperationRecord")
	authorityMenuApi := new(v1.OperationRecordApi)

	operationRecordRouter.Post("createSysOperationRecord", authorityMenuApi.CreateSysOperationRecord)             // 新建SysOperationRecord
	operationRecordRouter.Delete("deleteSysOperationRecord/:id", authorityMenuApi.DeleteSysOperationRecord)       // 删除SysOperationRecord
	operationRecordRouter.Delete("deleteSysOperationRecordByIds", authorityMenuApi.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
	operationRecordRouter.Get("findSysOperationRecord/:id", authorityMenuApi.FindSysOperationRecord)              // 根据ID获取SysOperationRecord
	operationRecordRouter.Get("getSysOperationRecordList", authorityMenuApi.GetSysOperationRecordList)            // 获取SysOperationRecord列表
}
