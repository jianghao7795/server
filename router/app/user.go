package app

import (
	v1 "server/api/v1/app"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserRouter struct{}

// InitUserRouter 初始化 User 路由信息
func (s *UserRouter) InitUserRouter(Router fiber.Router) {
	userRouter := Router.Group("frontend-user")
	var userApi = new(v1.UserApi)

	userRouter.Post("createUser", middleware.OperationRecord, userApi.CreateUser)             // 新建User
	userRouter.Delete("deleteUser/:id", middleware.OperationRecord, userApi.DeleteUser)       // 删除User
	userRouter.Delete("deleteUserByIds", middleware.OperationRecord, userApi.DeleteUserByIds) // 批量删除User
	userRouter.Put("updateUser/:id", middleware.OperationRecord, userApi.UpdateUser)          // 更新User

	userRouter.Get("findUser/:id", userApi.FindUser)   // 根据ID获取User
	userRouter.Get("getUserList", userApi.GetUserList) // 获取User列表

}
