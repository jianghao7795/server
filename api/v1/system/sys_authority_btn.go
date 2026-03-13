package system

import (
	global "server/model"
	"server/model/common/response"
	"server/model/system/request"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct{}

// @Tags AuthorityBtn
// @Summary 获取权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysAuthorityBtnReq true "菜单id, 角色id, 选中的按钮id"
// @Success 200 {object} response.Response{data=[]system.SysAuthorityBtn,msg=string} "返回列表成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c *fiber.Ctx) error {
	var req request.SysAuthorityBtnReq
	_ = c.QueryParser(&req)
	if res, err := authorityBtnService.GetAuthorityBtn(req); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败", c)
	} else {
		return response.OkWithDetailed(res, "查询成功", c)
	}
}

// @Tags AuthorityBtn
// @Summary 设置权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysAuthorityBtnReq true "菜单id, 角色id, 选中的按钮id"
// @Success 200 {object} response.Response{msg=string} "返回列表成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /authorityBtn/setAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c *fiber.Ctx) error {
	var req request.SysAuthorityBtnReq
	_ = c.BodyParser(&req)
	if err := authorityBtnService.SetAuthorityBtn(req); err != nil {
		global.LOG.Error("分配失败!", zap.Error(err))
		return response.FailWithMessage("分配失败", c)
	} else {
		return response.OkWithMessage("分配成功", c)
	}
}

// @Tags AuthorityBtn
// @Summary 设置权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /authorityBtn/canRemoveAuthorityBtn [post]
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := authorityBtnService.CanRemoveAuthorityBtn(id); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}
