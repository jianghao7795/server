package response

import (
	"github.com/gofiber/fiber/v3"
)

type ResponseMobile struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	ERROR401 = fiber.StatusUnauthorized
)

// 返回401 错误信息 data 和 string message信息返回
func FailWithDetailed401(data any, message string, c fiber.Ctx) error {
	return Result(ERROR401, data, message, c)
}
