package utils

import (
	"fmt"
	"server-fiber/model"
	"server-fiber/model/common/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler provides common error handling utilities
type ErrorHandler struct{}

// HandleAPIError handles common API errors with consistent logging and response
func (eh *ErrorHandler) HandleAPIError(c *fiber.Ctx, operation string, err error) error {
	if err == nil {
		return nil
	}

	// Log the error with context
	model.LOG.Error(fmt.Sprintf("%s failed", operation), zap.Error(err))

	// Return standardized error response
	return response.FailWithDetailed(map[string]string{
		"error": err.Error(),
	}, fmt.Sprintf("%s失败", operation), c)
}

// HandleValidationError handles parameter validation errors
func (eh *ErrorHandler) HandleValidationError(c *fiber.Ctx, field string, err error) error {
	model.LOG.Error(fmt.Sprintf("Validation failed for field %s", field), zap.Error(err))
	return response.FailWithMessage(fmt.Sprintf("参数验证失败: %s", err.Error()), c)
}

// HandleDatabaseError handles database operation errors
func (eh *ErrorHandler) HandleDatabaseError(c *fiber.Ctx, operation string, err error) error {
	model.LOG.Error(fmt.Sprintf("Database %s failed", operation), zap.Error(err))
	return response.FailWithMessage(fmt.Sprintf("数据库操作失败: %s", operation), c)
}

// HandleNotFoundError handles resource not found errors
func (eh *ErrorHandler) HandleNotFoundError(c *fiber.Ctx, resource string) error {
	model.LOG.Warn(fmt.Sprintf("%s not found", resource))
	return response.FailWithMessage(fmt.Sprintf("%s不存在", resource), c)
}

// Global error handler instance
var ErrorHandlerInstance = &ErrorHandler{}
