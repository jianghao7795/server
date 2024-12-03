package system

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"server-fiber/config"
	"server-fiber/model/system/request"
	"server-fiber/utils"

	"github.com/gookit/color"

	global "server-fiber/model"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

// WriteConfig mysql回写配置
func (h *MysqlInitHandler) WriteConfig(ctx context.Context) error {
	c, ok := ctx.Value("config").(config.Mysql)
	if !ok {
		return errors.New("mysql config invalid")
	}
	global.CONFIG.System.DbType = "mysql"
	global.CONFIG.Mysql = c
	cs := utils.StructToMap(global.CONFIG)
	for k, v := range cs {
		global.VIP.Set(k, v)
	}
	global.VIP.Set("jwt.signing-key", uuid.New().String())
	return global.VIP.WriteConfig()
}

// EnsureDB 创建数据库并初始化 mysql
func (h *MysqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (next context.Context, err error) {
	if s, ok := ctx.Value("dbtype").(string); !ok || s != "mysql" {
		return ctx, ErrDBTypeMismatch
	}
	dsn := conf.MysqlEmptyDsn()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.DBName)
	if err = createDatabase(dsn, "mysql", createSql); err != nil {
		return nil, err
	} // 创建数据库

	c := conf.ToMysqlConfig()
	type myString string
	next = context.WithValue(ctx, myString("config"), c)
	if c.Dbname == "" {
		return ctx, nil
	} // 如果没有数据库名, 则跳出初始化数据
	var db *gorm.DB
	if db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: true,    // 根据版本自动配置
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return ctx, err
	}
	global.CONFIG.AutoCode.Root, _ = filepath.Abs("..")

	next = context.WithValue(next, myString("db"), db)
	return next, err
}

func (h *MysqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	return createTables(ctx, inits)
}

func (h MysqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer func(c func()) { c() }(cancel)
	for _, init := range inits {
		if init.DataInserted(next) {
			color.Info.Printf(InitDataExist, Mysql, init.InitializerName())
			continue
		}
		if n, err := init.InitializeData(next); err != nil {
			color.Info.Printf(InitDataFailed, Mysql, Mysql, err)
			return err
		} else {
			next = n
			color.Info.Printf(InitDataSuccess, Mysql, init.InitializerName())
		}
	}
	color.Info.Printf(InitSuccess, Mysql)
	return nil
}
