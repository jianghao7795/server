package middleware

import (
	"server/model/common/response"
	"server/service/frontend"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func JWTAuthMiddleware(c fiber.Ctx) error {
	// 解决访问文件的401问题
	if strings.Contains(c.Get("Accept"), "image/") {
		code := c.Response().StatusCode()
		return c.Status(code).SendFile(strings.Join(strings.Split(c.Path(), "/")[2:], "/"))
	}

	authHeader := c.Get("Authorization")
	token := strings.Replace(authHeader, "Bearer ", "", 1)
	if token == "" {
		return response.FailWithMessage401("token 失效， 请重新登录", c)
	}
	_, err := frontend.ParseToken(token)
	// c.Locals("frontend_user", user)
	if err != nil {
		return response.FailWithMessage("token 失效， 请重新登录", c)
	}

	return c.Next() //
}
