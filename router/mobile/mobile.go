package mobile

import (
	fileUpload "server/api/v1/app"
	v1 "server/api/v1/mobile"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type MobileLoginRouter struct{}

func (m *MobileUserRouter) InitMobileLoginRouter(Router fiber.Router) {
	mobileLoginRouter := Router.Group("").Use(middleware.JWTAuthMobileMiddleware())
	var mobileLoginApi = new(v1.LoginApi)
	var registerApi = new(v1.RegisterMobile)
	{
		mobileLoginRouter.Post("login", mobileLoginApi.Login)
		mobileLoginRouter.Post("register", registerApi.Register)
	}
	mobileGetUserApi := Router.Group("").Use(middleware.JWTAuthMobileMiddleware())
	{
		mobileGetUserApi.Get("getUserInfo", mobileLoginApi.GetUserInfo)
		mobileGetUserApi.Put("updateUser", mobileLoginApi.UpdateMobileUser)
		mobileGetUserApi.Put("updatePassword", mobileLoginApi.UpdatePassword)
	}
	exaFileUploadAndDownloadApi := new(fileUpload.FileUploadAndDownloadApi)
	{
		mobileGetUserApi.Post("uploadImage", exaFileUploadAndDownloadApi.UploadFile)
	}

}
