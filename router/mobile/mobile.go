package mobile

import (
	fileUpload "server/api/v1/app"
	v1 "server/api/v1/mobile"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type MobileLoginRouter struct{}

func (m *MobileLoginRouter) InitMobileLoginRouter(Router fiber.Router) {
	mobileLoginRouter := Router.Group("mobile")
	var mobileLoginApi = new(v1.LoginApi)
	var registerApi = new(v1.RegisterMobile)
	{
		mobileLoginRouter.Post("login", mobileLoginApi.Login)
		mobileLoginRouter.Post("register", registerApi.Register)
	}
	mobileGetUserApi := mobileLoginRouter.Use(middleware.JWTAuthMobileMiddleware())
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
