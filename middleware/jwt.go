package middleware

import (
	"strings"

	global "server/model"
	"server/model/common/response"
	systemReq "server/model/system/request"
	systemService "server/service/system"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
)

var jwtService = new(systemService.JwtService)

// JWTAuth 使用 Fiber 内置 JWT 中间件（contrib/v3/jwt），RS256 + 自定义 Claims + 黑名单
func JWTAuth(c fiber.Ctx) error {
	return jwtMiddleware(c)
}

var jwtMiddleware = jwtware.New(jwtware.Config{
	// 跳过静态文件路径，避免 401
	Next: func(c fiber.Ctx) bool {
		return strings.Contains(c.Path(), "uploads/excel/") || strings.Contains(c.Path(), "uploads/file/")
	},
	SigningKey: jwtware.SigningKey{
		JWTAlg: jwtware.RS256,
		Key:    global.RunCONFIG.JWT.PublicKey,
	},
	Claims: &systemReq.CustomClaims{},
	ErrorHandler: func(c fiber.Ctx, err error) error {
		msg := "未登录或非法访问"
		if err != nil && strings.Contains(strings.ToLower(err.Error()), "expired") {
			msg = "授权已过期"
		} else if err != nil && err.Error() != "" {
			msg = err.Error()
		}
		return response.FailWithDetailed401(fiber.Map{"reload": true}, msg, c)
	},
	SuccessHandler: func(c fiber.Ctx) error {
		rawToken := c.Get("Authorization", "")
		tokenStr := strings.TrimPrefix(strings.TrimSpace(rawToken), "Bearer ")
		if tokenStr == "" {
			return response.FailWithDetailed401(fiber.Map{"reload": true}, "未登录或非法访问", c)
		}
		if jwtService.IsBlacklist(tokenStr) {
			return response.FailWithDetailed401(fiber.Map{"reload": true}, "您的帐户异地登陆或令牌失效", c)
		}
		token := jwtware.FromContext(c)
		if token != nil && token.Valid {
			if claims, ok := token.Claims.(*systemReq.CustomClaims); ok {
				c.Locals("claims", claims)
			}
		}
		return c.Next()
	},
})
