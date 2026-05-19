package response

import (
	global "server/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// Response 结构体
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	ERROR             = fiber.StatusBadRequest   // 错误返回的code 数据
	SUCCESS           = fiber.StatusOK           // 成功返回的code
	ERRORNotFound     = fiber.StatusNotFound     // 404错误
	ERRORUnauthorized = fiber.StatusUnauthorized // 401错误
	ERROR403          = fiber.StatusForbidden    // 403错误
)

// 底层的返回结果
func Result(code int, data any, msg string, c fiber.Ctx) error {
	// 返回的最终结果
	return c.Status(code).JSON(Response{
		code,
		data,
		msg,
	})
}

// 成功返回
func Ok(c fiber.Ctx) error {
	logMethodMessage("操作成功", nil, 2)
	return Result(SUCCESS, map[string]any{}, "操作成功", c)
}

// 成功返回 并带string信息返回
func OkWithMessage(message string, c fiber.Ctx) error {
	logMethodMessage("操作成功", nil, 2)
	return Result(SUCCESS, map[string]any{}, message, c)
}

// 成功返回 并带id信息返回
func OkWithId(message string, id uint, c fiber.Ctx) error {
	logMethodMessage("操作成功", nil, 2)
	return Result(SUCCESS, map[string]uint{
		"id": id,
	}, message, c)
}

// 成功返回 并带data信息返回
func OkWithData(data any, c fiber.Ctx) error {
	logMethodMessage("操作成功", nil, 2)
	return Result(SUCCESS, data, "操作成功", c)
}

// 成功返回 并带data 和 string message信息返回
func OkWithDetailed(data any, message string, c fiber.Ctx) error {
	logMethodMessage("操作成功", nil, 2)
	return Result(SUCCESS, data, message, c)
}

func logMethodMessage(msg string, err error, status int, fields ...zap.Field) { // status 1.warn 2.info 3.error，4.debug， 5.fatal
	fs := fields
	var logMethod func(string, ...zap.Field)
	switch status {
	case 1:
		logMethod = global.LOG.Warn
	case 2:
		logMethod = global.LOG.Info
	case 3:
		logMethod = global.LOG.Error
	case 4:
		logMethod = global.LOG.Debug
	case 5:
		logMethod = global.LOG.Fatal
	}
	if err != nil {
		fs = append([]zap.Field{zap.Error(err)}, fs...)
	}

	if len(fs) > 0 {
		logMethod(msg, fs...)
	} else {
		logMethod(msg)
	}
}

// 失败返回
func Fail(status int, err error, c fiber.Ctx) error {
	logMethodMessage("操作失败", err, status)
	return Result(ERROR, map[string]any{}, "操作失败", c)
}

func FailWithMessage(message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Error(err))
	return Result(ERROR, map[string]any{
		"msg": err.Error(),
	}, message, c)
}

func FailWithDetailed(data any, message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Any("data", data))
	return Result(ERROR, data, message, c)
}

// 返回400 错误信息 带上data信息
func FailWithDetailed400(data any, message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Any("data", data))
	return Result(ERROR, data, message, c)
}

// 返回400 错误信息 带上message信息
func FailWithMessage400(message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status)
	return Result(ERROR, map[string]any{}, message, c)
}

func FailWithMessage404(message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Int("code", ERRORNotFound))
	return Result(ERRORNotFound, map[string]any{}, message, c)
}

// 返回401 错误信息 带上message信息
func FailWithMessage401(message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Int("code", ERRORUnauthorized))
	return Result(ERRORUnauthorized, map[string]any{}, message, c)
}

// 返回403 错误信息 带上message信息
func FailWithMessage403(message string, status int, err error, c fiber.Ctx) error {
	logMethodMessage(message, err, status, zap.Int("code", ERROR403))
	return Result(ERROR403, map[string]any{}, message, c)
}
