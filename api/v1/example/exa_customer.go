package example

import (
	"server-fiber/global"
	"server-fiber/model/common/response"
	"server-fiber/model/example"
	"server-fiber/model/example/request"
	exampleRes "server-fiber/model/example/response"
	"server-fiber/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// @Tags ExaCustomer
// @Summary 创建客户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body example.ExaCustomer true "客户用户名, 客户手机号码"
// @Success 200 {object} response.Response{msg=string} "创建客户"
// @Router /customer/customer [post]
func (e *CustomerApi) CreateExaCustomer(c *fiber.Ctx) error {
	var customer example.ExaCustomer
	if err := c.BodyParser(&customer); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}

	if err := utils.Verify(customer, utils.CustomerVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	customer.SysUserID = utils.GetUserID(c)
	customer.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	if err := customerService.CreateExaCustomer(&customer); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithId("创建成功", customer.ID, c)
	}
}

// @Tags ExaCustomer
// @Summary 删除客户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path number true "客户ID"
// @Success 200 {object} response.Response{msg=string} "删除客户"
// @Router /customer/customer/:id [delete]
func (e *CustomerApi) DeleteExaCustomer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if id == 0 {
		return response.FailWithMessage("id不存在", c)
	}
	if err := customerService.DeleteExaCustomer(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// @Tags ExaCustomer
// @Summary 更新客户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path number true "客户ID"
// @Param data body example.ExaCustomer true "客户ID, 客户信息"
// @Success 200 {object} response.Response{msg=string,code=number} "更新客户信息"
// @Router /customer/customer/:id [put]
func (e *CustomerApi) UpdateExaCustomer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if id == 0 {
		return response.FailWithMessage("id不存在", c)
	}
	var customer example.ExaCustomer
	_ = c.BodyParser(&customer)
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
		return response.FailWithMessage("更新失败", c)
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
// @Success 200 {object} response.Response{data=exampleRes.ExaCustomerResponse,msg=string} "获取单一客户信息,返回包括客户详情"
// @Router /customer/customer/:id [get]
func (e *CustomerApi) GetExaCustomer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
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
// @Router /customer/customerList [get]
func (e *CustomerApi) GetExaCustomerList(c *fiber.Ctx) error {
	var pageInfo request.SearchCustomerParams
	_ = c.QueryParser(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	customerList, total, err := customerService.GetCustomerInfoList(utils.GetUserAuthorityId(c), &pageInfo)
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
