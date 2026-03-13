package mobile

import (
	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/mobile"
	mobileReq "server/model/mobile/request"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserApi struct{}

// CreateMobileUser 创建移动端用户
// @Tags MobileUser
// @Summary 创建移动端用户
// @Description 创建新的移动端用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body mobile.MobileUser true "移动端用户信息"
// @Success 200 {object} response.Response{msg=string} "创建移动端用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/createMobileUser [post]
func (userApi *UserApi) CreateMobileUser(c *fiber.Ctx) error {
	var mobileUser mobile.MobileUser
	if err := c.BodyParser(&mobileUser); err != nil {
		global.LOG.Error("获取用户数据失败", zap.Error(err))
		return response.FailWithMessage("获取用户数据失败", c)
	}
	if err := userService.CreateMobileUser(&mobileUser); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithMessage("创建成功", c)
	}
}

// DeleteMobileUser 删除移动端用户
// @Tags MobileUser
// @Summary 删除移动端用户
// @Description 根据用户ID删除指定移动端用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "用户ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除移动端用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/deleteMobileUser/{id} [delete]
func (userApi *UserApi) DeleteMobileUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := userService.DeleteMobileUser(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// DeleteMobileUserByIds 批量删除移动端用户
// @Tags MobileUser
// @Summary 批量删除移动端用户
// @Description 根据ID列表批量删除移动端用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "用户ID列表"
// @Success 200 {object} response.Response{msg=string} "创建移动端用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/deleteMobileUserByIds [delete]
func (userApi *UserApi) DeleteMobileUserByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	if err := c.BodyParser(&IDS); err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("批量获取id失败", c)
	}
	if err := userService.DeleteMobileUserByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithMessage("批量删除失败", c)
	} else {
		return response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateMobileUser 更新移动端用户
// @Tags MobileUser
// @Summary 更新移动端用户
// @Description 更新移动端用户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body mobile.MobileUser true "移动端用户信息"
// @Success 200 {object} response.Response{msg=string} "更新移动端用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/updateMobileUser [put]
func (userApi *UserApi) UpdateMobileUser(c *fiber.Ctx) error {
	var mobileUser mobile.MobileUser
	if err := c.BodyParser(&mobileUser); err != nil {
		global.LOG.Error("获取用户信息失败", zap.Error(err))
		return response.FailWithMessage("获取用户信息失败", c)
	}
	if err := userService.UpdateMobileUser(&mobileUser); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// FindMobileUser 用id查询MobileUser
// @Tags MobileUser
// @Summary 用id查询MobileUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query mobile.MobileUser true "用id查询MobileUser"
// @Success 200 {string} string "{"success":true,"data":mobile.MobileUser,"msg":"查询成功"}"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/findMobileUser/:id [get]
func (userApi *UserApi) FindMobileUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if mobileUser, err := userService.GetMobileUser(uint(id)); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败", c)
	} else {
		return response.OkWithData(mobileUser, c)
	}
}

// GetMobileUserList 分页获取MobileUser列表
// @Tags MobileUser
// @Summary 分页获取MobileUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query mobileReq.MobileUserSearch true "分页获取MobileUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobileUser/getMobileUserList [get]
func (userApi *UserApi) GetMobileUserList(c *fiber.Ctx) error {
	var pageInfo mobileReq.MobileUserSearch
	_ = c.QueryParser(&pageInfo)
	if list, total, err := userService.GetMobileUserInfoList(&pageInfo); err != nil {
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
