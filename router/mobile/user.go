package mobile

import (
	v1 "server/api/v1/mobile"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type MobileUserRouter struct{}

func (m *MobileUserRouter) InitMobileRouter(Router fiber.Router) {
	mobileUserRouterWithoutRecord := Router.Group("mobile")
	var mobileUserApi = new(v1.UserApi)

	mobileUserRouterWithoutRecord.Post("createMobileUser", middleware.OperationRecord, mobileUserApi.CreateMobileUser)             // 新建MobileUser
	mobileUserRouterWithoutRecord.Delete("deleteMobileUser/:id", middleware.OperationRecord, mobileUserApi.DeleteMobileUser)       // 删除MobileUser
	mobileUserRouterWithoutRecord.Delete("deleteMobileUserByIds", middleware.OperationRecord, mobileUserApi.DeleteMobileUserByIds) // 批量删除MobileUser
	mobileUserRouterWithoutRecord.Put("updateMobileUser/:id", middleware.OperationRecord, mobileUserApi.UpdateMobileUser)          // 更新MobileUser

	mobileUserRouterWithoutRecord.Get("findMobileUser/:id", mobileUserApi.FindMobileUser)   // 根据ID获取MobileUser
	mobileUserRouterWithoutRecord.Get("getMobileUserList", mobileUserApi.GetMobileUserList) // 获取MobileUser列表
}
