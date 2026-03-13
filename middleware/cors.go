package middleware

import (
	"server/config"
	global "server/model"
	"server/model/common/response"

	"github.com/gofiber/fiber/v3"
)

// Cors 直接放行所有跨域请求并放行所有 OPTIONS 方法
func Cors(c fiber.Ctx) error {
	method := c.Method()
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		return c.Next()
	}
	origin := c.Get("access-control-allow-origin")
	if origin == "" {
		c.Set("Access-Control-Allow-Origin", "*")
	} else {
		c.Set("Access-Control-Allow-Origin", origin)
	}
	c.Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
	c.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
	c.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Set("Access-Control-Allow-Credentials", "true")

	// 处理请求
	return c.Next()
}

// CorsByRules 按照配置处理跨域请求
func CorsByRules(c fiber.Ctx) error {
	// 放行全部
	if global.CONFIG.Cors.Mode == "allow-all" {
		return Cors(c)
	}
	whitelist := checkCors(c.Get("access-control-allow-origin"))
	// 通过检查, 添加请求头
	if whitelist != nil {
		c.Set("Access-Control-Allow-Origin", whitelist.AllowOrigin)
		c.Set("Access-Control-Allow-Headers", whitelist.AllowHeaders)
		c.Set("Access-Control-Allow-Methods", whitelist.AllowMethods)
		c.Set("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
		if whitelist.AllowCredentials {
			c.Set("Access-Control-Allow-Credentials", "true")
		}
	}

	// 严格白名单模式且未通过检查，直接拒绝处理请求
	if whitelist == nil && global.CONFIG.Cors.Mode == "strict-whitelist" {
		return response.FailWithMessage403("域名拒绝访问", c)
	} else {
		// 非严格白名单模式，无论是否通过检查均放行所有 OPTIONS 方法
		if c.Method() == "OPTIONS" {
			c.Status(fiber.StatusOK)
		}
	}

	// 处理请求
	return c.Next()
}

func checkCors(currentOrigin string) *config.CORSWhitelist {
	for _, whitelist := range global.CONFIG.Cors.Whitelist {
		// 遍历配置中的跨域头，寻找匹配项
		if currentOrigin == whitelist.AllowOrigin {
			return &whitelist
		}
	}
	return nil
}
