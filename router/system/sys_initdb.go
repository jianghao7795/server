package system

import (
	v1 "server/api/v1/system"

	"github.com/gofiber/fiber/v3"
)

type InitRouter struct{}

func (s *InitRouter) InitInitRouter(Router fiber.Router) {
	initRouter := Router.Group("init")
	dbApi := new(v1.DBApi)

	initRouter.Post("initdb", dbApi.InitDB)  // 创建Api
	initRouter.Get("checkdb", dbApi.CheckDB) // 检查是否初始化
}
