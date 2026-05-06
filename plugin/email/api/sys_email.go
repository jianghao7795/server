package api

import (
	"server/model/common/response"
	email_response "server/plugin/email/model/response"
	"server/plugin/email/service"

	"github.com/gofiber/fiber/v3"
)

type EmailApi struct{}

// @Tags System
// @Summary 发送测试邮件
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /email/emailTest [post]
func (s *EmailApi) EmailTest(c fiber.Ctx) error {
	if err := service.ServiceGroupApp.EmailTest(); err != nil {
		return response.FailWithMessage("发送失败", 3, err, c)
	} else {
		return response.OkWithData("发送成功", c)
	}
}

// @Tags System
// @Summary 发送邮件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body email_response.Email true "发送邮件必须的参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /email/sendEmail [post]
func (s *EmailApi) SendEmail(c fiber.Ctx) error {
	var email email_response.Email
	_ = c.Bind().Body(&email)
	if err := service.ServiceGroupApp.SendEmail(email.To, email.Subject, email.Body); err != nil {
		return response.FailWithMessage("发送失败", 3, err, c)
	} else {
		return response.OkWithData("发送成功", c)
	}
}
