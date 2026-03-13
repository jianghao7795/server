package response

import (
	"github.com/gofiber/fiber/v3"
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
	return Result(SUCCESS, map[string]any{}, "操作成功", c)
}

// 成功返回 并带string信息返回
func OkWithMessage(message string, c fiber.Ctx) error {
	return Result(SUCCESS, map[string]any{}, message, c)
}

// 成功返回 并带id信息返回
func OkWithId(message string, id uint, c fiber.Ctx) error {
	return Result(SUCCESS, map[string]uint{
		"id": id,
	}, message, c)
}

// 成功返回 并带data信息返回
func OkWithData(data any, c fiber.Ctx) error {
	return Result(SUCCESS, data, "操作成功", c)
}

// 成功返回 并带data 和 string message信息返回
func OkWithDetailed(data any, message string, c fiber.Ctx) error {
	return Result(SUCCESS, data, message, c)
}

// 失败返回
func Fail(c fiber.Ctx) error {
	return Result(ERROR, map[string]any{}, "操作失败", c)
}

func FailWithMessage(message string, c fiber.Ctx) error {
	return Result(ERROR, map[string]any{}, message, c)
}

func FailWithDetailed(data any, message string, c fiber.Ctx) error {
	return Result(ERROR, data, message, c)
}

// 返回400 错误信息 带上data信息
func FailWithDetailed400(data any, message string, c fiber.Ctx) error {
	return Result(ERROR, data, message, c)
}

// 返回400 错误信息 带上message信息
func FailWithMessage400(message string, c fiber.Ctx) error {
	return Result(ERROR, map[string]any{}, message, c)
}

func FailWithMessage404(message string, c fiber.Ctx) error {
	return Result(ERRORNotFound, map[string]any{}, message, c)
}

// 返回401 错误信息 带上message信息
func FailWithMessage401(message string, c fiber.Ctx) error {
	return Result(ERRORUnauthorized, map[string]any{}, message, c)
}

// 返回403 错误信息 带上message信息
func FailWithMessage403(message string, c fiber.Ctx) error {
	return Result(ERROR403, map[string]any{}, message, c)
}
