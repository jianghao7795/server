package plugin

import "github.com/gofiber/fiber/v3"

const (
	OnlyFuncName = "Plugin"
)

// Plugin 插件模式接口化
type Plugin interface {
	// Register 注册路由
	Register(group fiber.Router)

	// RouterPath 用户返回注册路由
	RouterPath() string
}
