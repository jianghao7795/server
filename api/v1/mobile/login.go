package mobile

import (
	"errors"
	"strconv"

	global "server/model"
	"server/model/common/response"
	"server/model/mobile"
	"server/model/mobile/request"
	"server/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type LoginApi struct{}

// Login 移动端用户登录
// @Tags Mobile Login
// @Summary 移动端用户登录
// @Description 移动端用户登录获取用户信息
// @Accept application/json
// @Produce application/json
// @Param data body mobile.Login true "登录信息"
// @Success 200 {object} response.Response{msg=string} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobile/login [post]
func (*LoginApi) Login(c fiber.Ctx) error {
	var l mobile.Login
	if err := c.Bind().Body(&l); err != nil {
		global.LOG.Error("获取登录数据失败", zap.Error(err))
		return response.FailWithMessage("获取登录数据失败", c)
	}
	if err := utils.Verify(l, utils.MobileLoginVerify); err != nil { // 验证用户密码的规则
		return response.FailWithMessage(err.Error(), c)
	}
	loginResponse, err := loginService.Login(&l)
	if err != nil {
		global.LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		return response.FailWithMessage400("用户名不存在或者密码错误", c)
	} else {
		return response.OkWithDetailed(loginResponse, "登录成功", c)
	}

}

// GetUserInfo 获取移动端用户信息
// @Tags Mobile Login
// @Summary 获取移动端用户信息
// @Description 根据用户ID获取移动端用户详细信息
// @Produce application/json
// @Param user_id header string true "用户ID"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobile/getUserInfo [get]
func (*LoginApi) GetUserInfo(c fiber.Ctx) error {
	authorization := c.Get("user_id", "")
	if authorization == "" {
		global.LOG.Error("获取user_id失败!", zap.Error(errors.New("失败")))
		return response.FailWithMessage400("获取失败", c)
	}
	authorityId, _ := strconv.Atoi(authorization)
	if user, err := loginService.GetUserInfo(uint(authorityId)); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage400("获取失败", c)
	} else {
		return response.OkWithDetailed(user, "获取成功", c)
	}

}

// UpdateMobileUser 更新移动端用户信息
// @Tags Mobile Login
// @Summary 更新移动端用户信息
// @Description 更新移动端用户的基本信息
// @Accept application/json
// @Produce application/json
// @Param user_id header string true "用户ID"
// @Param data body request.MobileUpdate true "用户更新信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobile/updateMobileUser [put]
func (*LoginApi) UpdateMobileUser(c fiber.Ctx) error {
	var data request.MobileUpdate
	if err := c.Bind().Body(&data); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	authorization := c.Get("user_id", "") // user_id 在请求头信息中

	if authorization == "" {
		global.LOG.Error("获取User_id 失败")
		return response.FailWithDetailed400(fiber.Map{"id": 0}, "更新失败", c)
	} else {
		authId, _ := strconv.Atoi(authorization)
		if authId == 0 {
			global.LOG.Error("user_id 转化错误")
		}
		if err := loginService.UpdateUser(&data, uint(authId)); err != nil {
			global.LOG.Error("更新用户失败!", zap.Error(err))
			return response.FailWithMessage("更新用户失败", c)
		} else {
			return response.OkWithDetailed(data, "更新成功", c)
		}
	}

}

// UpdatePassword 更新移动端用户密码
// @Tags Mobile Login
// @Summary 更新移动端用户密码
// @Description 更新移动端用户的登录密码
// @Accept application/json
// @Produce application/json
// @Param data body request.MobileUpdatePassword true "密码更新信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobile/updatePassword [put]
func (*LoginApi) UpdatePassword(c fiber.Ctx) error {
	var data request.MobileUpdatePassword
	if err := c.Bind().Body(&data); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}

	if err := utils.Verify(data, utils.MobileUpdatePasswordVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}

	if err := loginService.UpdatePassword(data); err != nil {
		global.LOG.Error("更新密码失败!", zap.Error(err))
		return response.FailWithMessage("更新用户密码失败", c)
	} else {
		return response.OkWithDetailed(data.NewPassword, "更新成功", c)
	}
}
