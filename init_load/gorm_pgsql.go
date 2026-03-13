package init_load

import (
	"errors"
	"server/config"
	"server/init_load/gorm_log"
	global "server/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormPgSql 初始化 Postgresql 数据库
func GormPgSql() (*gorm.DB, error) {
	p := global.CONFIG.Pgsql
	if p.Dbname == "" {
		return nil, errors.New("没设置数据库名")
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), gorm_log.Gorm.Config()); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		return db, nil
	}
}

// GormPgSqlByConfig 初始化 Postgresql 数据库 通过参数
func GormPgSqlByConfig(p config.DB) *gorm.DB {
	if p.Dbname == "" {
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), gorm_log.Gorm.Config()); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		return db
	}
}
