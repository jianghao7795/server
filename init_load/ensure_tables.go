package init_load

import (
	"context"

	"server/model/app"
	"server/model/example"
	sysModel "server/model/system"
	"server/service/system"

	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// 初始化数据库

const initOrderEnsureTables = system.InitOrderExternal - 1

type ensureTables struct{}

// auto run 初始化数据库
func init() {
	system.RegisterInit(initOrderEnsureTables, &ensureTables{})
}

func (ensureTables) InitializerName() string {
	return "ensure_tables_created"
}
func (e *ensureTables) InitializeData(ctx context.Context) (next context.Context, err error) {
	return ctx, nil
}

func (e *ensureTables) DataInserted(ctx context.Context) bool {
	return true
}

// 迁移数据
func (e *ensureTables) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	tables := []any{
		// sysModel.SysApi{},
		// sysModel.SysUser{},
		// sysModel.SysBaseMenu{},
		// sysModel.SysAuthority{},
		// sysModel.JwtBlacklist{},
		// sysModel.SysDictionary{},
		// sysModel.SysAutoCodeHistory{},
		// sysModel.SysOperationRecord{},
		// sysModel.SysDictionaryDetail{},
		// sysModel.SysBaseMenuParameter{},
		// sysModel.SysBaseMenuBtn{},
		// sysModel.SysAuthorityBtn{},
		// sysModel.SysAutoCode{},

		// adapter.CasbinRule{},

		// example.ExaFile{},
		// example.ExaCustomer{},
		// example.ExaFileChunk{},
		// example.ExaFileUploadAndDownload{},
	}
	for _, t := range tables {
		_ = db.AutoMigrate(&t)
		// 视图 authority_menu 会被当成表来创建，引发冲突错误（更新版本的gorm似乎不会）
		// 由于 AutoMigrate() 基本无需考虑错误，因此显式忽略
	}
	return ctx, nil
}

// 判断是否已经创建了表
func (e *ensureTables) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
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
		app.Praise{},
		app.User{},
	}
	yes := true
	for _, t := range tables {
		yes = yes && db.Migrator().HasTable(t)
	}
	return yes
}
