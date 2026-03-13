package system

import (
	global "server/model"
	"server/model/common/response"
	"server/model/system"
	"server/model/system/request"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DictionaryDetailApi struct{}

// @Tags SysDictionaryDetail
// @Summary 创建SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysDictionaryDetail true "SysDictionaryDetail模型"
// @Success 200 {object} response.Response{msg=string} "创建SysDictionaryDetail"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysDictionaryDetail/createSysDictionaryDetail [post]
func (s *DictionaryDetailApi) CreateSysDictionaryDetail(c *fiber.Ctx) error {
	var detail system.SysDictionaryDetail
	_ = c.BodyParser(&detail)
	if err := dictionaryDetailService.CreateSysDictionaryDetail(detail); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithMessage("创建成功", c)
	}
}

// @Tags SysDictionaryDetail
// @Summary 删除SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysDictionaryDetail true "SysDictionaryDetail模型"
// @Success 200 {object} response.Response{msg=string} "删除SysDictionaryDetail"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysDictionaryDetail/deleteSysDictionaryDetail/:id [delete]
func (s *DictionaryDetailApi) DeleteSysDictionaryDetail(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := dictionaryDetailService.DeleteSysDictionaryDetail(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// @Tags SysDictionaryDetail
// @Summary 更新SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysDictionaryDetail true "更新SysDictionaryDetail"
// @Success 200 {object} response.Response{msg=string} "更新SysDictionaryDetail"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysDictionaryDetail/updateSysDictionaryDetail [put]
func (s *DictionaryDetailApi) UpdateSysDictionaryDetail(c *fiber.Ctx) error {
	var detail system.SysDictionaryDetail
	_ = c.BodyParser(&detail)
	if err := dictionaryDetailService.UpdateSysDictionaryDetail(&detail); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// @Tags SysDictionaryDetail
// @Summary 用id查询SysDictionaryDetail
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query system.SysDictionaryDetail true "用id查询SysDictionaryDetail"
// @Success 200 {object} response.Response{data=system.SysDictionaryDetail,msg=string} "用id查询SysDictionaryDetail"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysDictionaryDetail/findSysDictionaryDetail/:id [get]
func (s *DictionaryDetailApi) FindSysDictionaryDetail(c *fiber.Ctx) error {
	var detail system.SysDictionaryDetail
	id, _ := c.ParamsInt("id")
	detail.ID = uint(id)
	if err := utils.Verify(detail, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if resysDictionaryDetail, err := dictionaryDetailService.GetSysDictionaryDetail(detail.ID); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败", c)
	} else {
		return response.OkWithDetailed(resysDictionaryDetail, "查询成功", c)
	}
}

// @Tags SysDictionaryDetail
// @Summary 分页获取SysDictionaryDetail列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysDictionaryDetailSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]system.SysDictionaryDetail,total=int64,page=int,pageSize=int},msg=string} "分页获取SysDictionaryDetail列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /sysDictionaryDetail/getSysDictionaryDetailList [get]
func (s *DictionaryDetailApi) GetSysDictionaryDetailList(c *fiber.Ctx) error {
	var pageInfo request.SysDictionaryDetailSearch
	_ = c.QueryParser(&pageInfo)
	if list, total, err := dictionaryDetailService.GetSysDictionaryDetailInfoList(pageInfo); err != nil {
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
