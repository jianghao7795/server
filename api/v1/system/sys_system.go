package system

import (
	global "server/model"
	"server/model/common/response"
	"server/model/system"
	systemRes "server/model/system/response"
	"server/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type SystemApi struct{}

// @Tags System
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{data=systemRes.SysConfigResponse,msg=string} "获取配置文件内容,返回包括系统配置"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /system/getSystemConfig [get]
func (s *SystemApi) GetSystemConfig(c fiber.Ctx) error {
	if config, err := systemConfigService.GetSystemConfig(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(systemRes.SysConfigResponse{Config: config}, "获取成功", c)
	}
}

// @Tags System
// @Summary 设置配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body system.System true "设置配置文件内容"
// @Success 200 {object} response.Response{msg=string} "设置配置文件内容"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /system/setSystemConfig [put]
func (s *SystemApi) SetSystemConfig(c fiber.Ctx) error {
	var sys system.System
	if err := c.Bind().Body(&sys); err != nil {
		global.LOG.Error("获取配置数据失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	if err := systemConfigService.SetSystemConfig(sys); err != nil {
		global.LOG.Error("设置失败!", zap.Error(err))
		return response.FailWithMessage("设置失败", c)
	} else {
		return response.OkWithData("设置成功", c)
	}
}

// @Tags System
// @Summary 重启系统
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{msg=string} "重启系统"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /system/reloadSystem [post]
func (s *SystemApi) ReloadSystem(c fiber.Ctx) error {
	err := utils.Reload()
	if err != nil {
		global.LOG.Error("重启系统失败!", zap.Error(err))
		return response.FailWithMessage("重启系统失败", c)
	} else {
		return response.OkWithMessage("重启系统成功", c)
	}
}

// @Tags System
// @Summary 获取服务器信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{data=systemRes.ServerInfo,msg=string} "获取服务器信息"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /system/getServerInfo [post]
func (s *SystemApi) GetServerInfo(c fiber.Ctx) error {
	if server, err := systemConfigService.GetServerInfo(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(server, "获取成功", c)
	}
}
