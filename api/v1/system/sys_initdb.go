package system

import (
	global "server/model"
	"server/model/app"
	"server/model/common/response"
	"server/model/example"
	sysModel "server/model/system"
	"server/model/system/request"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
)

type DBApi struct{}

// InitDB
// @Tags InitDB
// @Summary 初始化用户数据库
// @Produce  application/json
// @Param data body request.InitDB true "初始化数据库参数"
// @Success 200 {object} response.Response{msg=string} "初始化用户数据库"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /init/initdb [post]
func (i *DBApi) InitDB(c fiber.Ctx) error {
	if global.DB != nil {
		return response.FailWithMessage("已存在数据库配置", 3, nil, c)
	}
	var dbInfo request.InitDB
	if err := c.Bind().Query(&dbInfo); err != nil {
		return response.FailWithMessage("参数校验不通过", 3, err, c)
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		return response.FailWithMessage("自动创建数据库失败，请查看后台日志，检查后在进行初始化", 3, err, c)
	}
	return response.OkWithData("自动创建数据库成功", c)
}

// CheckDB
// @Tags CheckDB
// @Summary 是否进行初始化
// @Produce  application/json
// @Success 200 {object} response.Response{msg=string} "初始化用户数据库"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /init/checkdb [get]
func (i *DBApi) CheckDB(c fiber.Ctx) error {
	var (
		message  = ""
		needInit = true
	)

	if global.DB == nil {
		message = "数据库连接失败"
		needInit = false
	}
	if i.hasTable() {
		message = "数据库初始化成功"
		return response.OkWithDetailed(needInit, message, c)
	} else {
		message = "数据库初始化失败：请查看后台日志"
		needInit = true
	}
	return response.FailWithDetailed400(needInit, message, 3, nil, c)
}

// hasTable检查数据库中是否存在表
func (initDB *DBApi) hasTable() bool {
	tables := []any{
		sysModel.SysApi{},
		sysModel.SysUser{},
		sysModel.SysBaseMenu{},
		sysModel.SysAuthority{},
		sysModel.JwtBlacklist{},
		sysModel.SysDictionary{},
		sysModel.SysAutoCodeHistory{},
		sysModel.SysOperationRecord{},
		sysModel.SysDictionaryDetail{},
		sysModel.SysBaseMenuParameter{},
		sysModel.SysBaseMenuBtn{},
		sysModel.SysAuthorityBtn{},
		sysModel.SysAutoCode{},

		adapter.CasbinRule{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		app.Article{},
		app.ArticleTag{},
		app.Tag{},
		app.BaseMessage{},
		app.Comment{},
		app.Ip{},
		global.Praise{},
		app.User{},
	}
	yes := true
	for _, t := range tables {
		yes = yes && global.DB.Migrator().HasTable(t)
	}
	return yes
}
