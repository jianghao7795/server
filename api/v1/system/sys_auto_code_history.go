package system

import (
	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	systemReq "server/model/system/request"
	systemRes "server/model/system/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AutoCodeHistoryApi struct{}

// First
// @Tags AutoCode
// @Summary 获取meta信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "请求参数"
// @Success 200 {object} response.Response{data=string,msg=string} "获取meta信息"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /autoCode/getMeta [get]
func (a *AutoCodeHistoryApi) First(c fiber.Ctx) error {
	var info request.GetById
	_ = c.Bind().Query(&info)
	data, err := autoCodeHistoryService.First(&info)
	if err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	return response.OkWithDetailed(data, "获取成功", c)
}

// Delete
// @Tags AutoCode
// @Summary 删除回滚记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "请求参数"
// @Success 200 {object} response.Response{msg=string} "删除回滚记录"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /autoCode/delSysHistory [get]
func (a *AutoCodeHistoryApi) Delete(c fiber.Ctx) error {
	var info request.GetById
	_ = c.Bind().Query(&info)
	err := autoCodeHistoryService.Delete(&info)
	if err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	}
	return response.OkWithMessage("删除成功", c)
}

// RollBack
// @Tags AutoCode
// @Summary 回滚自动生成代码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.RollBack true "请求参数"
// @Success 200 {object} response.Response{msg=string} "回滚自动生成代码"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /autoCode/rollback [get]
func (a *AutoCodeHistoryApi) RollBack(c fiber.Ctx) error {
	var info systemReq.RollBack
	_ = c.Bind().Query(&info)
	if err := autoCodeHistoryService.RollBack(&info); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	return response.OkWithMessage("回滚成功", c)
}

// GetList
// @Tags AutoCode
// @Summary 查询回滚记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SysAutoHistory true "请求参数"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]systemRes.AutoCodeHistory,total=int64,page=int,pageSize=int},msg=string} "查询回滚记录,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /autoCode/getSysHistory [post]
func (a *AutoCodeHistoryApi) GetList(c fiber.Ctx) error {
	var search systemReq.SysAutoHistory
	_ = c.Bind().Query(&search)
	var list []systemRes.AutoCodeHistory
	var total int64
	var err error
	list, total, err = autoCodeHistoryService.GetList(search.PageInfo)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	}
	return response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}
