package system

import (
	"strconv"

	global "server/model"
	"server/model/common/response"
	"server/model/system/request"
	modelSystemResponse "server/model/system/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct{}

// @Tags AuthorityBtn
// @Summary 获取权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysAuthorityBtnReq true "菜单id, 角色id, 选中的按钮id"
// @Success 200 {object} response.Response{data=modelSystemResponse.SysAuthorityBtnRes,msg=string,code=integer} "返回列表成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c fiber.Ctx) error {
	var req request.SysAuthorityBtnReq
	_ = c.Bind().Query(&req)
	var res modelSystemResponse.SysAuthorityBtnRes
	var err error
	res, err = authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
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
// @Success 200 {object} response.Response{msg=string,data=modelSystemResponse.SysAuthorityBtnRes,code=integer} "返回列表成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /authorityBtn/setAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c fiber.Ctx) error {
	var req request.SysAuthorityBtnReq
	_ = c.Bind().Body(&req)
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
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := authorityBtnService.CanRemoveAuthorityBtn(id); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}
