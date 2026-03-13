package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type JwtRouter struct{}

func (s *JwtRouter) InitJwtRouter(Router fiber.Router) {
	jwtRouter := Router.Group("jwt")
	jwtApi := new(v1.JwtApi)

	jwtRouter.Post("jsonInBlacklist", middleware.OperationRecord, jwtApi.JsonInBlacklist) // jwt加入黑名单

}
