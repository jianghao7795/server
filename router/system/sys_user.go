package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router fiber.Router) {
	userRouter := Router.Group("user")
	baseApi := new(v1.BaseApi)

	userRouter.Post("admin_register", middleware.OperationRecord, baseApi.Register)               // 管理员注册账号
	userRouter.Post("changePassword", middleware.OperationRecord, baseApi.ChangePassword)         // 用户修改密码
	userRouter.Post("setUserAuthority", middleware.OperationRecord, baseApi.SetUserAuthority)     // 设置用户权限
	userRouter.Delete("deleteUser/:id", middleware.OperationRecord, baseApi.DeleteUser)           // 删除用户
	userRouter.Put("setUserInfo", middleware.OperationRecord, baseApi.SetUserInfo)                // 设置用户信息
	userRouter.Put("setSelfInfo", middleware.OperationRecord, baseApi.SetSelfInfo)                // 设置自身信息
	userRouter.Post("setUserAuthorities", middleware.OperationRecord, baseApi.SetUserAuthorities) // 设置用户权限组
	userRouter.Post("resetPassword", middleware.OperationRecord, baseApi.ResetPassword)           // 重置密码

	userRouter.Get("getUserList", baseApi.GetUserList)   // 分页获取用户列表
	userRouter.Get("getUserInfo", baseApi.GetUserInfo)   // 获取自身信息
	userRouter.Get("getUserCount", baseApi.GetUserCount) // 获取用户数
	userRouter.Get("getFlow", baseApi.GetFlow)           // 获取摸个ip流量
}
