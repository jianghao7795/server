package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitAuthorityRouter(Router fiber.Router) {
	authorityRouter := Router.Group("authority")
	authorityApi := new(v1.AuthorityApi)

	authorityRouter.Post("createAuthority", middleware.OperationRecord, authorityApi.CreateAuthority)   // 创建角色
	authorityRouter.Delete("DeleteAuthority", middleware.OperationRecord, authorityApi.DeleteAuthority) // 删除角色
	authorityRouter.Put("updateAuthority", middleware.OperationRecord, authorityApi.UpdateAuthority)    // 更新角色
	authorityRouter.Post("copyAuthority", middleware.OperationRecord, authorityApi.CopyAuthority)       // 拷贝角色
	authorityRouter.Post("setDataAuthority", middleware.OperationRecord, authorityApi.SetDataAuthority) // 设置角色资源权限

	authorityRouter.Get("getAuthorityList", authorityApi.GetAuthorityList) // 获取角色列表

}
