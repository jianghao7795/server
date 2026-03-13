package system

import (
	v1 "server/api/v1/system"

	"github.com/gofiber/fiber/v2"
)

type AutoCodeHistoryRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeHistoryRouter(Router fiber.Router) {
	autoCodeHistoryRouter := Router.Group("autoCode")
	autoCodeHistoryApi := new(v1.AutoCodeHistoryApi)

	autoCodeHistoryRouter.Post("getMeta", autoCodeHistoryApi.First)         // 根据id获取meta信息
	autoCodeHistoryRouter.Post("rollback", autoCodeHistoryApi.RollBack)     // 回滚
	autoCodeHistoryRouter.Post("delSysHistory", autoCodeHistoryApi.Delete)  // 删除回滚记录
	autoCodeHistoryRouter.Post("getSysHistory", autoCodeHistoryApi.GetList) // 获取回滚记录分页

}
