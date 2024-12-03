package system

import (
	global "server-fiber/model"
	"server-fiber/model/common/request"
	"server-fiber/model/common/response"
	systemReq "server-fiber/model/system/request"

	"github.com/gofiber/fiber/v2"
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
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取meta信息"
// @Router /autoCode/getMeta [get]
func (a *AutoCodeHistoryApi) First(c *fiber.Ctx) error {
	var info request.GetById
	_ = c.QueryParser(&info)
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
// @Router /autoCode/delSysHistory [get]
func (a *AutoCodeHistoryApi) Delete(c *fiber.Ctx) error {
	var info request.GetById
	_ = c.QueryParser(&info)
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
// @Router /autoCode/rollback [get]
func (a *AutoCodeHistoryApi) RollBack(c *fiber.Ctx) error {
	var info systemReq.RollBack
	_ = c.QueryParser(&info)
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
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "查询回滚记录,返回包括列表,总数,页码,每页数量"
// @Router /autoCode/getSysHistory [post]
func (a *AutoCodeHistoryApi) GetList(c *fiber.Ctx) error {
	var search systemReq.SysAutoHistory
	_ = c.QueryParser(&search)
	list, total, err := autoCodeHistoryService.GetList(search.PageInfo)
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
