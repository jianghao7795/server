package ws

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// wsPlugin 结构体定义了一个WebSocket插件，包含日志输出对象和buffer大小
type wsPlugin struct {
	logger    *zap.Logger // 日志输出对象
	manageBuf int64       // buffer大小
	// registeredMsgHandler map[int32]func(biz.IMessage) bool // 消息处理函数映射
	// checkMap             map[string]biz.CheckFunc          // 用户校验函数映射

	// admin     biz.IManage
	// adminCase *biz.AdminCase
}

// Register 方法注册了WebSocket路由
func (w *wsPlugin) Register(g fiber.Router) {
	// ws 为身份校验函数
	g.Get("/ws", func(c fiber.Ctx) error {
		// 实现WebSocket连接的逻辑
		return nil
	})
	g.Post("/sendMsg", func(c fiber.Ctx) error {
		// 实现发送消息的逻辑
		return nil
	})
}

// RouterPath 方法返回插件的路由路径
func (w *wsPlugin) RouterPath() string {
	return "ws"
}

// 假设有一个全局变量来存储插件实例，以便在其他地方使用
var globalWSPlugin *wsPlugin

func init() {
	// 初始化插件实例，这里可以添加插件初始化逻辑
	globalWSPlugin = &wsPlugin{
		logger:    zap.NewExample(), // 示例日志初始化，实际应用中应替换为合适的日志配置
		manageBuf: 1024,
	}
	// 注册插件，这里假设有一个全局的Fiber应用实例
	app := fiber.New()
	globalWSPlugin.Register(app)
	// 启动服务器等后续操作...
}
