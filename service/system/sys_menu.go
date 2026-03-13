package system

import (
	"errors"
	global "server/model"
	"server/model/common/request"
	"server/model/system"
	"strconv"

	"gorm.io/gorm"
)

//@function: getMenuTreeMap
//@description: 获取路由总树map
//@param: authorityId string
//@return: err error, treeMap map[string][]model.SysMenu

func (menuService *MenuService) getMenuTreeMap(authorityId string) (treeMap map[string][]system.SysMenu, err error) {
	var allMenus []system.SysMenu
	var btns []system.SysAuthorityBtn
	treeMap = make(map[string][]system.SysMenu)
	err = global.DB.Where("authority_id = ?", authorityId).Order("sort").Preload("Parameters").Find(&allMenus).Error
	if err != nil {
		return
	}
	err = global.DB.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
	if err != nil {
		return
	}
	btnMap := make(map[uint]map[string]string)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]string)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}
	for _, v := range allMenus {
		v.Btns = btnMap[v.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: GetMenuTree
//@description: 获取动态菜单树
//@param: authorityId string
//@return: err error, menus []model.SysMenu

func (menuService *MenuService) GetMenuTree(authorityId string) (menus []system.SysMenu, err error) {
	menuTree, err := menuService.getMenuTreeMap(authorityId)
	menus = menuTree["0"]
	for i := range len(menus) {
		err = menuService.getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}

//@function: getChildrenList
//@description: 获取子菜单
//@param: menu *model.SysMenu, treeMap map[string][]model.SysMenu
//@return: err error

func (menuService *MenuService) getChildrenList(menu *system.SysMenu, treeMap map[string][]system.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := range len(menu.Children) {
		err = menuService.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@function: GetInfoList
//@description: 获取路由分页
//@return: err error, list any, total int64

func (menuService *MenuService) GetInfoList() (list []system.SysBaseMenu, total int64, err error) {
	var menuList []system.SysBaseMenu
	treeMap, err := menuService.getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := range len(menuList) {
		err = menuService.getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, total, err
}

//@function: getBaseChildrenList
//@description: 获取菜单的子菜单
//@param: menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu
//@return: err error

func (menuService *MenuService) getBaseChildrenList(menu *system.SysBaseMenu, treeMap map[string][]system.SysBaseMenu) (err error) {
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	for i := range len(menu.Children) {
		err = menuService.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@function: AddBaseMenu
//@description: 添加基础路由
//@param: menu model.SysBaseMenu
//@return: error

func (menuService *MenuService) AddBaseMenu(menu system.SysBaseMenu) error {
	if !errors.Is(global.DB.Where("name = ?", menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.DB.Create(&menu).Error
}

//@function: getBaseMenuTreeMap
//@description: 获取路由总树map
//@return: err error, treeMap map[string][]model.SysBaseMenu

func (menuService *MenuService) getBaseMenuTreeMap() (treeMap map[string][]system.SysBaseMenu, err error) {
	var allMenus []system.SysBaseMenu
	treeMap = make(map[string][]system.SysBaseMenu)
	err = global.DB.Order("sort").Preload("MenuBtn").Preload("Parameters").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: GetBaseMenuTree
//@description: 获取基础路由树
//@return: err error, menus []model.SysBaseMenu

func (menuService *MenuService) GetBaseMenuTree() (menus []system.SysBaseMenu, err error) {
	treeMap, err := menuService.getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := range len(menus) {
		err = menuService.getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

//@function: AddMenuAuthority
//@description: 为角色增加menu树
//@param: menus []model.SysBaseMenu, authorityId string
//@return: err error

func (menuService *MenuService) AddMenuAuthority(menus []system.SysBaseMenu, authorityId string) (err error) {
	var auth system.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus
	err = AuthorityServiceApp.SetMenuAuthority(&auth)
	return err
}

//@function: GetMenuAuthority
//@description: 查看当前角色树
//@param: info *request.GetAuthorityId
//@return: err error, menus []model.SysMenu

func (menuService *MenuService) GetMenuAuthority(info *request.GetAuthorityId) (menus []system.SysMenu, err error) {
	err = global.DB.Where("authority_id = ? ", info.AuthorityId).Order("sort").Find(&menus).Error
	// sql := "SELECT authority_menu.keep_alive,authority_menu.default_menu,authority_menu.created_at,authority_menu.updated_at,authority_menu.deleted_at,authority_menu.menu_level,authority_menu.parent_id,authority_menu.path,authority_menu.`name`,authority_menu.hidden,authority_menu.component,authority_menu.title,authority_menu.icon,authority_menu.sort,authority_menu.menu_id,authority_menu.authority_id FROM authority_menu WHERE authority_menu.authority_id = ? ORDER BY authority_menu.sort ASC"
	// err = global.DB.Raw(sql, authorityId).Scan(&menus).Error
	return menus, err
}
