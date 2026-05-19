package response

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
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
func FailWithDetailed401(data any, message string, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, 3, zap.Any("data", data), zap.Int("code", ERROR401))
	return Result(ERROR401, data, message, c)
}
