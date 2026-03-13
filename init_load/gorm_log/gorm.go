package gorm_log

import (
	"log"
	"os"
	global "server/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBBASE interface {
	GetLogMode() string
}

var Gorm = new(gormDB)

type gormDB struct{}

// Config 记录 gorm sql语句
// gorm 自定义配置
// 注意 QueryFields 模式会根据当前 model 的所有字段名称进行 select
// DisableForeignKeyConstraintWhenMigrating 在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true
func (g *gormDB) Config() *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, QueryFields: true, SkipDefaultTransaction: global.CONFIG.Mysql.SkipTransaction, PrepareStmt: true}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      false,
	})
	var logMode DBBASE
	switch global.CONFIG.System.DbType {
	case "mysql":
		logMode = &global.CONFIG.Mysql
	case "pgsql":
		logMode = &global.CONFIG.Pgsql
	case "sqlite":
		logMode = &global.CONFIG.Sqlite
	default:
		logMode = &global.CONFIG.Mysql
	}

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
