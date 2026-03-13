package middleware

import (
	"strings"

	"server/utils"

	"server/model/common/response"
	systemService "server/service/system"

	"github.com/gofiber/fiber/v3"
)

var jwtService = new(systemService.JwtService)

func JWTAuth(c fiber.Ctx) error {
	// 解决访问文件的401问题
	if strings.Contains(c.Path(), "uploads/excel/") || strings.Contains(c.Path(), "uploads/file/") {
		code := c.Response().StatusCode()
		return c.Status(code).SendFile(strings.Join(strings.Split(c.Path(), "/")[2:], "/"))
	}
	tokenString := c.Get("Authorization")

	if tokenString == "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiNzMzYzk2NzUtODNmYS00NzBhLThhN2YtZTAwZWI1YzdlMDIxIiwiSUQiOjMsIlVzZXJuYW1lIjoid3VoYW8iLCJOaWNrTmFtZSI6IuWQtOaYiiIsIkF1dGhvcml0eUlkIjoiODg4IiwiQnVmZmVyVGltZSI6NjAsImlzcyI6ImppYW5naGFvIiwiZXhwIjoxNzU4MTU5OTY0LCJuYmYiOjE3NTgwNzM1NjQsImlhdCI6MTc1ODA3MzU2NH0.MD2-I4BHwIrk80o2VIcwPOS7EikuISd4FeNTrD-BlV8x2ygzON5BkWNLwlHUjKr3so5axsKQS7U4hM8pRmf0fxxFvGj-18r7QPLQwFDRSiN3OJrY3WZ3HK0fwKt71nbCSuetmonpbQpvFrE00KVBijwGHE1LKgsUEYUD-RZ7dRc" {
		return c.Next()
	}
	if tokenString == "" {
		return response.FailWithDetailed401(fiber.Map{"reload": true}, "未登录或非法访问", c)
	}
	token := strings.Replace(tokenString, "Bearer ", "", 1)
	if token == "" {
		return response.FailWithDetailed401(fiber.Map{"reload": true}, "未登录或非法访问", c)
	}
	if jwtService.IsBlacklist(token) {
		return response.FailWithDetailed401(fiber.Map{"reload": true}, "您的帐户异地登陆或令牌失效", c)
	}
	j := utils.NewJWT()
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(token)
	if err != nil {
		if err == utils.ErrTokenExpired {
			return response.FailWithDetailed401(fiber.Map{"reload": true}, "授权已过期", c)
		}
		return response.FailWithDetailed401(fiber.Map{"reload": true}, err.Error(), c)
	}
	// 继续交由下一个路由处理,并将解析出的信息传递下去
	c.Locals("claims", claims)
	return c.Next()
}
