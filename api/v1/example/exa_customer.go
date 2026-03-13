package example

import (
	"strconv"

	global "server/model"
	"server/model/common/response"
	"server/model/example"
	"server/model/example/request"
	exampleRes "server/model/example/response"
	"server/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// CreateExaCustomer 创建客户
// @Tags ExaCustomer
// @Summary 创建客户
// @Description 创建新的客户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body example.ExaCustomer true "客户信息"
// @Success 200 {object} response.Response{msg=string} "创建客户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /customer/customer [post]
func (e *CustomerApi) CreateExaCustomer(c fiber.Ctx) error {
	var customer example.ExaCustomer
	if err := c.Bind().Body(&customer); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}

	if err := utils.Verify(customer, utils.CustomerVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	var err error
	if customer.SysUserID, err = utils.GetUserID(c); err != nil {
		global.LOG.Warn("用户查询错误", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	customer.SysUserAuthorityID, err = utils.GetUserAuthorityId(c)
	if err != nil {
		global.LOG.Warn("获取用户权限失败", zap.Error(err))
		return response.FailWithMessage("新增用户失败"+err.Error(), c)
	}
	if err := customerService.CreateExaCustomer(&customer); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithId("创建成功", customer.ID, c)
	}
}

// DeleteExaCustomer 删除客户
// @Tags ExaCustomer
// @Summary 删除客户
// @Description 根据客户ID删除指定客户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "客户ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除客户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /customer/customer/{id} [delete]
func (e *CustomerApi) DeleteExaCustomer(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("请传id参数", c)
	}
	if id == 0 {
		return response.FailWithMessage("id传递错误", c)
	}
	if err := customerService.DeleteExaCustomer(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// UpdateExaCustomer 更新客户信息
// @Tags ExaCustomer
// @Summary 更新客户信息
// @Description 根据客户ID更新客户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "客户ID" minimum(1)
// @Param data body example.ExaCustomer true "客户信息"
// @Success 200 {object} response.Response{msg=string} "更新客户信息成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /customer/customer/{id} [put]
func (e *CustomerApi) UpdateExaCustomer(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return response.FailWithMessage("id不存在", c)
	}
	var customer example.ExaCustomer
	if err := c.Bind().Body(&customer); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if customer.ID != uint(id) {
		return response.FailWithMessage("数据不一致（id）", c)
	}
	if err := utils.Verify(customer, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := utils.Verify(customer, utils.CustomerVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := customerService.UpdateExaCustomer(&customer); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// @Tags ExaCustomer
// @Summary 获取单一客户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path number true "客户ID"
// @Success 200 {object} response.Response{data=object,msg=string} "获取单一客户信息,返回包括客户详情"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /customer/customer/:id [get]
func (e *CustomerApi) GetExaCustomer(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return response.FailWithMessage("id不存在", c)
	}
	data, err := customerService.GetExaCustomer(uint(id))
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(exampleRes.ExaCustomerResponse{Customer: data}, "获取成功", c)
	}
}

// @Tags ExaCustomer
// @Summary 分页获取权限客户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult{list=example.ExaCustomer[]},msg=string} "分页获取权限客户列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /customer/customerList [get]
func (e *CustomerApi) GetExaCustomerList(c fiber.Ctx) error {
	var pageInfo request.SearchCustomerParams
	_ = c.Bind().Query(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	authorityId, err := utils.GetUserAuthorityId(c)
	if err != nil {
		global.LOG.Error("获取用户权限Id失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	customerList, total, err := customerService.GetCustomerInfoList(authorityId, &pageInfo)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     customerList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
