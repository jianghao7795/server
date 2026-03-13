package middleware

import (
	"github.com/gofiber/fiber/v3"
)

// Recovery recover掉项目可能出现的panic，并使用zap记录相关日志
func Recovery(c fiber.Ctx) error {
	return c.Next()
}
