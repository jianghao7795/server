package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type DictionaryDetailRouter struct{}

func (s *DictionaryDetailRouter) InitSysDictionaryDetailRouter(Router fiber.Router) {
	dictionaryDetailRouter := Router.Group("sysDictionaryDetail")
	sysDictionaryDetailApi := new(v1.DictionaryDetailApi)

	dictionaryDetailRouter.Post("createSysDictionaryDetail", middleware.OperationRecord, sysDictionaryDetailApi.CreateSysDictionaryDetail)       // 新建SysDictionaryDetail
	dictionaryDetailRouter.Delete("deleteSysDictionaryDetail/:id", middleware.OperationRecord, sysDictionaryDetailApi.DeleteSysDictionaryDetail) // 删除SysDictionaryDetail
	dictionaryDetailRouter.Put("updateSysDictionaryDetail/:id", middleware.OperationRecord, sysDictionaryDetailApi.UpdateSysDictionaryDetail)    // 更新SysDictionaryDetail

	dictionaryDetailRouter.Get("findSysDictionaryDetail/:id", sysDictionaryDetailApi.FindSysDictionaryDetail)   // 根据ID获取SysDictionaryDetail
	dictionaryDetailRouter.Get("getSysDictionaryDetailList", sysDictionaryDetailApi.GetSysDictionaryDetailList) // 获取SysDictionaryDetail列表

}
