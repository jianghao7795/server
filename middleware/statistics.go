package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func Statistics(ctx fiber.Ctx) error {
	return ctx.Next()
}
