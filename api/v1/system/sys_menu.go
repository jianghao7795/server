package system

import (
	"strconv"

	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/system"
	systemReq "server/model/system/request"
	systemRes "server/model/system/response"
	"server/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AuthorityMenuApi struct{}

// GetMenu 获取用户动态路由
// @Tags AuthorityMenu
// @Summary 获取用户动态路由
// @Description 根据用户权限获取可访问的菜单树
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=[]system.SysMenu,msg=string} "获取用户动态路由成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/getMenu [get]
func (a *AuthorityMenuApi) GetMenu(c fiber.Ctx) error {
	authorityId, err := utils.GetUserAuthorityId(c)
	if err != nil {
		global.LOG.Error("获取权限id失败", zap.Error(err))
		return response.FailWithMessage("获取权限id", c)
	}
	if menus, err := menuService.GetMenuTree(authorityId); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		if menus == nil {
			menus = []system.SysMenu{}
		}
		return response.OkWithDetailed(menus, "获取成功", c)
	}
}

// GetBaseMenuTree 获取基础菜单树
// @Tags AuthorityMenu
// @Summary 获取基础菜单树
// @Description 获取系统所有基础菜单的树形结构
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=[]system.SysBaseMenu,msg=string} "获取基础菜单树成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/getBaseMenuTree [get]
func (a *AuthorityMenuApi) GetBaseMenuTree(c fiber.Ctx) error {
	if menus, err := menuService.GetBaseMenuTree(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(menus, "获取成功", c)
	}
}

// AddMenuAuthority 增加菜单和角色关联关系
// @Tags AuthorityMenu
// @Summary 增加菜单和角色关联关系
// @Description 为指定角色添加菜单权限
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.AddMenuAuthorityInfo true "菜单权限信息"
// @Success 200 {object} response.Response{msg=string} "增加菜单和角色关联关系成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/addMenuAuthority [post]
func (a *AuthorityMenuApi) AddMenuAuthority(c fiber.Ctx) error {
	var authorityMenu systemReq.AddMenuAuthorityInfo
	_ = c.Bind().Body(&authorityMenu)
	if err := utils.Verify(authorityMenu, utils.AuthorityIdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := menuService.AddMenuAuthority(authorityMenu.Menus, authorityMenu.AuthorityId); err != nil {
		global.LOG.Error("添加失败!", zap.Error(err))
		return response.FailWithMessage("添加失败", c)
	} else {
		return response.OkWithMessage("添加成功", c)
	}
}

// GetMenuAuthority 获取指定角色菜单
// @Tags AuthorityMenu
// @Summary 获取指定角色菜单
// @Description 获取指定角色可访问的菜单列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param authorityId query string true "角色ID"
// @Success 200 {object} response.Response{data=[]system.SysMenu,msg=string} "获取指定角色菜单成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/getMenuAuthority [get]
func (a *AuthorityMenuApi) GetMenuAuthority(c fiber.Ctx) error {
	var param request.GetAuthorityId
	_ = c.Bind().Query(&param)
	if err := utils.Verify(param, utils.AuthorityIdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if menus, err := menuService.GetMenuAuthority(&param); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "获取失败", c)
	} else {
		return response.OkWithDetailed(menus, "获取成功", c)
	}
}

// @Tags Menu
// @Summary 新增菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysBaseMenu true "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success 200 {object} response.Response{msg=string} "新增菜单"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/addBaseMenu [post]
func (a *AuthorityMenuApi) AddBaseMenu(c fiber.Ctx) error {
	var menu system.SysBaseMenu
	_ = c.Bind().Body(&menu)
	if err := utils.Verify(menu, utils.MenuVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := utils.Verify(menu.Meta, utils.MenuMetaVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := menuService.AddBaseMenu(menu); err != nil {
		global.LOG.Error("添加失败!", zap.Error(err))

		return response.FailWithMessage("添加失败", c)
	} else {
		return response.OkWithMessage("添加成功", c)
	}
}

// @Tags Menu
// @Summary 删除菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "菜单id"
// @Success 200 {object} response.Response{msg=string} "删除菜单"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/deleteBaseMenu [post]
func (a *AuthorityMenuApi) DeleteBaseMenu(c fiber.Ctx) error {
	var menu request.GetById
	menu.ID, _ = strconv.Atoi(c.Params("id"))
	if err := utils.Verify(menu, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := baseMenuService.DeleteBaseMenu(menu.ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// @Tags Menu
// @Summary 更新菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysBaseMenu true "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success 200 {object} response.Response{msg=string} "更新菜单"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/updateBaseMenu [post]
func (a *AuthorityMenuApi) UpdateBaseMenu(c fiber.Ctx) error {
	var menu system.SysBaseMenu
	_ = c.Bind().Body(&menu)
	if err := utils.Verify(menu, utils.MenuVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := utils.Verify(menu.Meta, utils.MenuMetaVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err := baseMenuService.UpdateBaseMenu(menu); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// @Tags Menu
// @Summary 根据id获取菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "菜单id"
// @Success 200 {object} response.Response{data=system.SysBaseMenu,msg=string} "根据id获取菜单,返回包括系统菜单列表"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/getBaseMenuById/:id [get]
func (a *AuthorityMenuApi) GetBaseMenuById(c fiber.Ctx) error {
	var idInfo request.GetById
	idInfo.ID, _ = strconv.Atoi(c.Params("id"))
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if menu, err := baseMenuService.GetBaseMenuById(idInfo.ID); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(menu, "获取成功", c)
	}
	// return response.OkWithDetailed("", "获取成功", c)
}

// @Tags Menu
// @Summary 分页获取基础menu列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]system.SysBaseMenu,total=int64,page=int,pageSize=int},msg=string} "分页获取基础menu列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /menu/getMenuList [post]
func (a *AuthorityMenuApi) GetMenuList(c fiber.Ctx) error {
	var pageInfo request.PageInfo
	_ = c.Bind().Query(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if menuList, total, err := menuService.GetInfoList(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     menuList,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
