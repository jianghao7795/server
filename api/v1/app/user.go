package app

import (
	"errors"
	global "server/model"
	"server/model/app"
	appReq "server/model/app/request"
	"server/model/common/request"
	"server/model/common/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CreateUser 创建用户
// @Tags User
// @Summary 创建用户
// @Description 创建新的应用用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.User true "用户信息"
// @Success 200 {object} response.Response{msg=string} "创建用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/createUser [post]
func (userApi *UserApi) CreateUser(c *fiber.Ctx) error {
	var user app.User
	err := c.BodyParser(&user)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := userService.CreateUser(&user); err != nil {
		global.LOG.Error(err.Error(), zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		return response.OkWithId("创建成功", user.ID, c)
	}
}

// DeleteUser 删除用户
// @Tags User
// @Summary 删除用户
// @Description 根据用户ID删除指定用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "用户ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /frontend-user/deleteUser/{id} [delete]
func (userApi *UserApi) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败!", zap.Error(err))
		return response.FailWithMessage("获取id失败，传正确的id", c)
	}
	if err := userService.DeleteUser(id); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// DeleteUserByIds 批量删除用户
// @Tags User
// @Summary 批量删除用户
// @Description 根据ID列表批量删除用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "用户ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/deleteUserByIds [delete]
func (userApi *UserApi) DeleteUserByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	err := c.BodyParser(&IDS)
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if err := userService.DeleteUserByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithMessage("批量删除失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateUser 更新用户
// @Tags User
// @Summary 更新用户
// @Description 根据用户ID更新用户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "用户ID" minimum(1)
// @Param data body app.User true "用户信息"
// @Success 200 {object} response.Response{msg=string} "更新用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/updateUser/{id} [put]
func (userApi *UserApi) UpdateUser(c *fiber.Ctx) error {
	var user app.User
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	var notFound bool
	notFound, err = userService.FindIsUser(id)
	if notFound && err != nil {
		global.LOG.Error("未找到，该用户!", zap.Error(errors.New("未找到，该用户")))
		return response.FailWithMessage("未找到，该用户", c)
	}
	err = c.BodyParser(&user)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败"+err.Error(), c)
	}
	user.ID = uint(id)
	if err := userService.UpdateUser(&user); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// FindUser 用id查询User
// @Tags User
// @Summary 用id查询User
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query app.User true "用id查询User"
// @Success 200 {object} response.Response{msg=string} "查询成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/findUser/:id [get]
func (userApi *UserApi) FindUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if user, err := userService.GetUser(uint(id)); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败"+err.Error(), c)
	} else {
		return response.OkWithData(user, c)
	}
}

// GetUserList 分页获取User列表
// @Tags User
// @Summary 分页获取User列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appReq.UserSearch true "分页获取User列表"
// @Success 200 {object} response.Response{msg=string,code=number,data=response.PageResult{list=app.User[]}} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/getUserList [get]
func (userApi *UserApi) GetUserList(c *fiber.Ctx) error {
	var pageInfo appReq.UserSearch
	_ = c.QueryParser(&pageInfo)
	if list, total, err := userService.GetUserInfoList(&pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
