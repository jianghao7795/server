package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthorityBtnRouter struct{}

func (s *AuthorityBtnRouter) InitAuthorityBtnRouterRouter(Router fiber.Router) {
	authorityRouter := Router.Group("authorityBtn")
	authorityBtnApi := new(v1.AuthorityBtnApi)

	authorityRouter.Post("setAuthorityBtn", middleware.OperationRecord, authorityBtnApi.SetAuthorityBtn)
	authorityRouter.Delete("canRemoveAuthorityBtn/:id", middleware.OperationRecord, authorityBtnApi.CanRemoveAuthorityBtn)

	authorityRouter.Get("getAuthorityBtn", authorityBtnApi.GetAuthorityBtn)

}
