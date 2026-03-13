package system

import (
	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/system"
	systemReq "server/model/system/request"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SystemApiApi struct{}

// @Tags SysApi
// @Summary 创建基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysApi true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "创建基础api"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/createApi [post]
func (s *SystemApiApi) CreateApi(c *fiber.Ctx) (err error) {
	var api system.SysApi
	err = c.BodyParser(&api)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = utils.Verify(api, utils.ApiVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err = apiService.CreateApi(&api); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{"msg": err.Error()}, "创建失败", c)
	} else {
		return response.OkWithId("创建成功", api.ID, c)
	}
}

// @Tags SysApi
// @Summary 删除api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path integer true "ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除api"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/delete/{id} [delete]
func (s *SystemApiApi) DeleteApi(c *fiber.Ctx) error {
	var api system.SysApi
	id, _ := c.ParamsInt("id")
	api.ID = uint(id)
	if err := utils.Verify(api.MODEL, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := apiService.DeleteApi(api); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// GetApiList @Tags SysApi
// @Summary 分页获取API列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SearchApiParams true "分页获取API列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/getApiList [get]
func (s *SystemApiApi) GetApiList(c *fiber.Ctx) error {
	var pageInfo systemReq.SearchApiParams
	_ = c.QueryParser(&pageInfo)
	// log.Println("pageInfo: ", pageInfo.ApiGroup)
	if list, total, err := apiService.GetAPIInfoList(&pageInfo); err != nil {
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

// GetApiById todo
// @Tags SysApi
// @Summary 根据id获取api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "根据id获取api"
// @Success 200 {object} response.Response{data=system.SysApi,msg=string} "根据id获取api,返回包括api详情"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/getApiById [post]
func (s *SystemApiApi) GetApiById(c *fiber.Ctx) error {
	var idInfo request.GetById
	idInfo.ID, _ = c.ParamsInt("id")
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	api, err := apiService.GetApiById(idInfo.ID)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithData(api, c)
	}
}

// UpdateApi @Tags SysApi
// @Summary 修改基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysApi true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "修改基础api"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/updateApi [put]
func (s *SystemApiApi) UpdateApi(c *fiber.Ctx) error {
	var api system.SysApi
	err := c.BodyParser(&api)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := utils.Verify(api, utils.ApiVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := apiService.UpdateApi(&api); err != nil {
		global.LOG.Error("修改失败!", zap.Error(err))
		return response.FailWithMessage("修改失败", c)
	} else {
		return response.OkWithMessage("修改成功", c)
	}
}

// @Tags SysApi
// @Summary 获取所有的Api 不分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]system.SysApi,msg=string} "获取所有的Api 不分页,返回包括api列表"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/getAllApis [get]
func (s *SystemApiApi) GetAllApis(c *fiber.Ctx) error {
	if apis, err := apiService.GetAllApis(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(apis, "获取成功", c)
	}
}

// @Tags SysApi
// @Summary 删除选中Api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {object} response.Response{msg=string} "删除选中Api"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /api/deleteApisByIds [delete]
func (s *SystemApiApi) DeleteApisByIds(c *fiber.Ctx) error {
	var ids request.IdsReq
	_ = c.QueryParser(&ids)
	if err := apiService.DeleteApisByIds(ids); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}
