package model

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// ModelSet Model 层依赖集合
// 主要用于提供全局变量的初始化
var ModelSet = wire.NewSet(
	ProvideGlobalVars,
)

// GlobalVars 全局变量结构
type GlobalVars struct {
	DB        *gorm.DB
	DBList    map[string]*gorm.DB
	REDIS     *redis.Client
	VIPER     *viper.Viper
	LOG       *zap.Logger
}

// ProvideGlobalVars 提供全局变量
func ProvideGlobalVars(
	db *gorm.DB,
	redis *redis.Client,
	vip *viper.Viper,
	logger *zap.Logger,
) *GlobalVars {
	// 设置全局变量
	DB = db
	REDIS = redis
	VIPER = vip
	LOG = logger
	
	return &GlobalVars{
		DB:     db,
		REDIS:  redis,
		VIPER:  vip,
		LOG:    logger,
	}
}
