package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router fiber.Router) {
	baseRouter := Router.Group("base")
	baseApi := new(v1.BaseApi)

	baseRouter.Post("login", baseApi.Login)
	baseRouter.Get("captcha", middleware.NeedInit, baseApi.Captcha)
	baseRouter.Get("captcha/img", baseApi.CaptchaImg)
	baseRouter.Post("getToken/login", baseApi.LoginToken)
}
