package middleware

import (
	"strings"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"

	global "server/model"
	"server/model/common/response"
	"server/utils"
)

func JWTAuthMobileMiddleware() fiber.Handler {
	if global.RunCONFIG.JWT.PublicKey == nil {
		return func(c fiber.Ctx) error {
			return response.FailWithMessage403("JWT 公钥未初始化", 3, nil, c)
		}
	}

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    global.RunCONFIG.JWT.PublicKey,
		},
		Claims: &utils.MobileClaims{},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return response.FailWithMessage401("token 失效，请重新登录", 3, err, c)
		},
		SuccessHandler: func(c fiber.Ctx) error {
			token := jwtware.FromContext(c)
			if token == nil || !token.Valid {
				return response.FailWithMessage401("token 失效，请重新登录", 3, nil, c)
			}

			claims, ok := token.Claims.(*utils.MobileClaims)
			if !ok || claims.ID == 0 {
				return response.FailWithMessage401("token 失效，请重新登录", 3, nil, c)
			}

			c.Locals("mobile_claims", claims)
			c.Locals("user_id", claims.ID)
			return c.Next()
		},
	})

	return func(c fiber.Ctx) error {
		if strings.Contains(c.Get("Accept"), "image/") {
			code := c.Response().StatusCode()
			return c.Status(code).SendFile(strings.Join(strings.Split(c.Path(), "/")[2:], "/"))
		}
		return jwtMiddleware(c)
	}
}
