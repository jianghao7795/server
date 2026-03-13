package router

import (
	"server/plugin/email/api"

	"github.com/gofiber/fiber/v3"
)

type EmailRouter struct{}

func (s *EmailRouter) InitEmailRouter(router fiber.Router) {
	EmailApi := api.ApiGroupApp.EmailApi.EmailTest
	SendEmail := api.ApiGroupApp.EmailApi.SendEmail
	{
		router.Post("emailTest", EmailApi)  // 发送测试邮件
		router.Post("sendEmail", SendEmail) // 发送邮件
	}
}
