package init_load

import (
	global "server/model"
	"server/plugin/email"
	"server/utils/plugin"

	"github.com/gofiber/fiber/v3"
)

func PluginInit(group fiber.Router, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
	}
}

func InstallPlugin(PrivateGroup fiber.Router) {
	//  添加跟角色挂钩权限的插件 示例 本地示例模式于在线仓库模式注意上方的import 可以自行切换 效果相同
	PluginInit(PrivateGroup, email.CreateEmailPlug(
		global.CONFIG.Email.To,
		global.CONFIG.Email.From,
		global.CONFIG.Email.Host,
		global.CONFIG.Email.Secret,
		global.CONFIG.Email.Nickname,
		global.CONFIG.Email.Port,
		global.CONFIG.Email.IsSSL,
	))
}
