package middleware

import "github.com/gofiber/fiber/v3"

func JWTAuthMiddleware(c fiber.Ctx) error {
	return JWTAuth(c)
}
