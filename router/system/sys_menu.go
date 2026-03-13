package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type MenuRouter struct{}

func (s *MenuRouter) InitMenuRouter(Router fiber.Router) {
	menuRouter := Router.Group("menu")
	authorityMenuApi := new(v1.AuthorityMenuApi)

	menuRouter.Post("addBaseMenu", middleware.OperationRecord, authorityMenuApi.AddBaseMenu)             // 新增菜单
	menuRouter.Post("addMenuAuthority", middleware.OperationRecord, authorityMenuApi.AddMenuAuthority)   //	增加menu和角色关联关系
	menuRouter.Delete("deleteBaseMenu/:id", middleware.OperationRecord, authorityMenuApi.DeleteBaseMenu) // 删除菜单
	menuRouter.Put("updateBaseMenu", middleware.OperationRecord, authorityMenuApi.UpdateBaseMenu)        // 更新菜单

	menuRouter.Get("getMenu", authorityMenuApi.GetMenu)                     // 获取菜单树
	menuRouter.Get("getMenuList", authorityMenuApi.GetMenuList)             // 分页获取基础menu列表
	menuRouter.Get("getBaseMenuTree", authorityMenuApi.GetBaseMenuTree)     // 获取用户动态路由
	menuRouter.Get("getMenuAuthority", authorityMenuApi.GetMenuAuthority)   // 获取指定角色menu
	menuRouter.Get("getBaseMenuById/:id", authorityMenuApi.GetBaseMenuById) // 根据id获取菜单

}
