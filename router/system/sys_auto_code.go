package system

import (
	v1 "server/api/v1/system"

	"github.com/gofiber/fiber/v2"
)

type AutoCodeRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeRouter(Router fiber.Router) {
	autoCodeRouter := Router.Group("autoCode")
	autoCodeApi := new(v1.AutoCodeApi)

	autoCodeRouter.Get("getDB", autoCodeApi.GetDB)                  // 获取数据库
	autoCodeRouter.Get("getTables", autoCodeApi.GetTables)          // 获取对应数据库的表
	autoCodeRouter.Get("getColumn", autoCodeApi.GetColumn)          // 获取指定表所有字段信息
	autoCodeRouter.Post("preview", autoCodeApi.PreviewTemp)         // 获取自动创建代码预览
	autoCodeRouter.Post("createTemp", autoCodeApi.CreateTemp)       // 创建自动化代码
	autoCodeRouter.Post("createPackage", autoCodeApi.CreatePackage) // 创建package包
	autoCodeRouter.Post("getPackage", autoCodeApi.GetPackage)       // 获取package包
	autoCodeRouter.Post("gelPackage", autoCodeApi.DelPackage)       //  删除package包

}
