package system

import (
	"strconv"

	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/system"
	systemReq "server/model/system/request"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type OperationRecordApi struct{}

// @Tags SysOperationRecord
// @Summary 创建SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysOperationRecord true "创建SysOperationRecord"
// @Success 200 {object} response.Response{msg=string} "创建SysOperationRecord"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysOperationRecord/createSysOperationRecord [post]
func (s *OperationRecordApi) CreateSysOperationRecord(c fiber.Ctx) error {
	var sysOperationRecord system.SysOperationRecord
	_ = c.Bind().Body(&sysOperationRecord)
	if err := operationRecordService.CreateSysOperationRecord(&sysOperationRecord); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithMessage("创建成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysOperationRecord true "SysOperationRecord模型"
// @Success 200 {object} response.Response{msg=string} "删除SysOperationRecord"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
func (s *OperationRecordApi) DeleteSysOperationRecord(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := operationRecordService.DeleteSysOperationRecord(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 批量删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SysOperationRecord"
// @Success 200 {object} response.Response{msg=string} "批量删除SysOperationRecord"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysOperationRecord/deleteSysOperationRecordByIds [delete]
func (s *OperationRecordApi) DeleteSysOperationRecordByIds(c fiber.Ctx) error {
	var IDS request.IdsReq
	_ = c.Bind().Body(&IDS)
	if err := operationRecordService.DeleteSysOperationRecordByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithMessage("批量删除失败", c)
	} else {
		return response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 用id查询SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query system.SysOperationRecord true "Id"
// @Success 200 {object} response.Response{data=system.SysOperationRecord,msg=string} "用id查询SysOperationRecord"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysOperationRecord/findSysOperationRecord/:id [get]
func (s *OperationRecordApi) FindSysOperationRecord(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if respSysOperationRecord, err := operationRecordService.GetSysOperationRecord(uint(id)); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败", c)
	} else {
		return response.OkWithDetailed(respSysOperationRecord, "查询成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 分页获取SysOperationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysOperationRecordSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]system.SysOperationRecord,total=int64,page=int,pageSize=int},msg=string} "分页获取SysOperationRecord列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysOperationRecord/getSysOperationRecordList [get]
func (s *OperationRecordApi) GetSysOperationRecordList(c fiber.Ctx) error {
	var pageInfo systemReq.SysOperationRecordSearch
	_ = c.Bind().Query(&pageInfo)
	if pageInfo.TypePort == system.Backend {
		if list, total, err := operationRecordService.GetSysOperationRecordInfoList(&pageInfo); err != nil {
			global.LOG.Error("获取失败!", zap.Error(err))
			return response.FailWithMessage("获取失败", c)
		} else {
			return response.OkWithDetailed(response.PageResult{
				List:     list,
				Total:    total,
				Page:     pageInfo.Page,
				PageSize: pageInfo.PageSize,
			}, "获取成功", c)
		}
	} else {
		if list, total, err := operationRecordService.GetSysOperationRecordInfoFrontendList(pageInfo); err != nil {
			global.LOG.Error("获取失败!", zap.Error(err))
			return response.FailWithMessage("获取失败", c)
		} else {
			return response.OkWithDetailed(response.PageResult{
				List:     list,
				Total:    total,
				Page:     pageInfo.Page,
				PageSize: pageInfo.PageSize,
			}, "获取成功", c)
		}
	}
}
