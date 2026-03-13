package example

import (
	v1 "server/api/v1/example"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type CustomerRouter struct{}

func (e *CustomerRouter) InitCustomerRouter(Router fiber.Router) {
	customerRouter := Router.Group("customer")
	exaCustomerApi := new(v1.CustomerApi)

	customerRouter.Post("customer", middleware.OperationRecord, exaCustomerApi.CreateExaCustomer)       // 创建客户
	customerRouter.Put("customer/:id", middleware.OperationRecord, exaCustomerApi.UpdateExaCustomer)    // 更新客户
	customerRouter.Delete("customer/:id", middleware.OperationRecord, exaCustomerApi.DeleteExaCustomer) // 删除客户

	customerRouter.Get("customer/:id", exaCustomerApi.GetExaCustomer)     // 获取单一客户信息
	customerRouter.Get("customerList", exaCustomerApi.GetExaCustomerList) // 获取客户列表
}
