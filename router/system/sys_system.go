package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type SysRouter struct{}

func (s *SysRouter) InitSystemRouter(Router fiber.Router) {
	sysRouter := Router.Group("system")
	systemApi := new(v1.SystemApi)

	sysRouter.Put("setSystemConfig", middleware.OperationRecord, systemApi.SetSystemConfig) // 设置配置文件内容
	sysRouter.Post("reloadSystem", middleware.OperationRecord, systemApi.ReloadSystem)      // 重启服务

	sysRouter.Get("getSystemConfig", systemApi.GetSystemConfig) // 获取配置文件内容
	sysRouter.Get("getServerInfo", systemApi.GetServerInfo)     // 获取服务器信息
}
