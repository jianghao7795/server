package middleware

import (
	"server-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

// ValidationMiddleware provides common validation utilities
type ValidationMiddleware struct{}

// ValidateBody validates request body and binds to struct
func (vm *ValidationMiddleware) ValidateBody(c *fiber.Ctx, dest interface{}) error {
	if err := c.BodyParser(dest); err != nil {
		return utils.ErrorHandlerInstance.HandleValidationError(c, "request body", err)
	}
	return nil
}

// ValidateParams validates URL parameters
func (vm *ValidationMiddleware) ValidateParams(c *fiber.Ctx, paramName string) (string, error) {
	value := c.Params(paramName)
	if value == "" {
		return "", utils.ErrorHandlerInstance.HandleValidationError(c, paramName, fiber.NewError(400, "missing required parameter"))
	}
	return value, nil
}

// ValidateQuery validates query parameters
func (vm *ValidationMiddleware) ValidateQuery(c *fiber.Ctx, queryName string) (string, error) {
	value := c.Query(queryName)
	if value == "" {
		return "", utils.ErrorHandlerInstance.HandleValidationError(c, queryName, fiber.NewError(400, "missing required query parameter"))
	}
	return value, nil
}

// ValidateID validates and converts ID parameter to uint
func (vm *ValidationMiddleware) ValidateID(c *fiber.Ctx) (uint, error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return 0, utils.ErrorHandlerInstance.HandleValidationError(c, "id", err)
	}
	if id <= 0 {
		return 0, utils.ErrorHandlerInstance.HandleValidationError(c, "id", fiber.NewError(400, "invalid ID"))
	}
	return uint(id), nil
}

// Global validation middleware instance
var ValidationMiddlewareInstance = &ValidationMiddleware{}
