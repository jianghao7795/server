package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router fiber.Router) {
	apiRouter := Router.Group("api")
	apiRouterApi := new(v1.SystemApiApi)

	apiRouter.Post("createApi", middleware.OperationRecord, apiRouterApi.CreateApi)               // 创建Api
	apiRouter.Delete("DeleteApi/:id", middleware.OperationRecord, apiRouterApi.DeleteApi)         // 删除Api
	apiRouter.Put("updateApi/:id", middleware.OperationRecord, apiRouterApi.UpdateApi)            // 更新api
	apiRouter.Delete("DeleteApisByIds", middleware.OperationRecord, apiRouterApi.DeleteApisByIds) // 删除选中api

	apiRouter.Get("getApiById/:id", apiRouterApi.GetApiById) // 获取单条Api消息
	apiRouter.Get("getAllApis", apiRouterApi.GetAllApis)     // 获取所有api
	apiRouter.Get("getApiList", apiRouterApi.GetApiList)     // 获取Api列表

}
