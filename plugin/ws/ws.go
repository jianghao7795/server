package ws

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type wsPlugin struct {
	logger    *zap.Logger
	manageBuf int64
}

func (w *wsPlugin) Register(g fiber.Router) {
	g.Get("/ws", func(c fiber.Ctx) error {
		return nil
	})
	g.Post("/sendMsg", func(c fiber.Ctx) error {
		return nil
	})
}

func (w *wsPlugin) RouterPath() string {
	return "ws"
}

var globalWSPlugin *wsPlugin

func init() {
	globalWSPlugin = &wsPlugin{
		logger:    zap.NewExample(),
		manageBuf: 1024,
	}
}
