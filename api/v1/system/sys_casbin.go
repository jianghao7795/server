package system

import (
	global "server/model"
	"server/model/common/response"
	"server/model/system/request"
	"server/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type CasbinApi struct{}

// UpdateCasbin 更新角色API权限
// @Tags Casbin
// @Summary 更新角色API权限
// @Description 更新指定角色的API权限配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限配置信息"
// @Success 200 {object} response.Response{msg=string} "更新角色API权限成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /casbin/UpdateCasbin [post]
func (cas *CasbinApi) UpdateCasbin(c fiber.Ctx) error {
	var cmr request.CasbinInReceive
	if err := c.Bind().Body(&cmr); err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := utils.Verify(cmr, utils.AuthorityIdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := casbinService.UpdateCasbin(cmr.AuthorityId, cmr.CasbinInfos); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// GetPolicyPathByAuthorityId 获取权限列表
// @Tags Casbin
// @Summary 获取权限列表
// @Description 根据角色ID获取权限策略列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response{data=fiber.Map{paths=[]request.CasbinInfo},msg=string} "获取权限列表成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /casbin/getPolicyPathByAuthorityId/{id} [get]
func (cas *CasbinApi) GetPolicyPathByAuthorityId(c fiber.Ctx) error {
	var casbin request.CasbinInReceive
	_ = c.Bind().Query(&casbin)
	casbin.AuthorityId = c.Params("id")
	if err := utils.Verify(casbin, utils.AuthorityIdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	paths := casbinService.GetPolicyPathByAuthorityId(casbin.AuthorityId)
	return response.OkWithDetailed(paths, "获取成功", c)
}
