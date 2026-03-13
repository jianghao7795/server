package init_load

import (
	"errors"
	"server/config"
	"server/init_load/gorm_log"
	global "server/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysql 初始化Mysql数据库

func GormMysql() (*gorm.DB, error) {
	m := global.CONFIG.Mysql
	if m.Dbname == "" {
		return nil, errors.New("no config database")
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
		// SkipDefaultTransaction:    m.SkipTransaction, // 禁用默认事务
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gorm_log.Gorm.Config()); err != nil {
		return nil, err
	} else {

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// db = db.Debug() // 线上注释
		if m.Debug {
			db = db.Debug()
		}
		return db, nil
	}
}

// GormMysqlByConfig 初始化Mysql数据库用过传入配置
func GormMysqlByConfig(m config.DB) (*gorm.DB, error) {
	if m.Dbname == "" {
		return nil, errors.New("请配置数据库")
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
		DisableDatetimePrecision:  true,    // 禁用datetime 精度
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gorm_log.Gorm.Config()); err != nil {
		return nil, err
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
