//go:build wireinject
// +build wireinject

package core

import (
	"server-fiber/init_load"
	"server-fiber/model"
	"server-fiber/router"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitializeApp 初始化应用
func InitializeApp() (*init_load.RouterApp, func(), error) {
	wire.Build(
		// 核心依赖
		ViperInit,
		ZapInit,
		init_load.ProvideGorm,
		init_load.ProvideRedis,
		init_load.ProvideTimer,

		// Router 层 - 使用 wire.NewSet 组织路由
		router.RouterSet,

		// Fiber App
		ProvideFiberApp,

		// Router 初始化
		init_load.ProvideRouter,
	)
	return nil, nil, nil
}

// ViperInit 初始化 Viper
func ViperInit() (*viper.Viper, error) {
	return viperInit()
}

// ZapInit 初始化 Zap
func ZapInit(vip *viper.Viper) (*zap.Logger, error) {
	if err := vip.Unmarshal(&model.CONFIG); err != nil {
		return nil, err
	}
	model.VIP = vip
	logger, err := zapInit()
	if err != nil {
		return nil, err
	}
	model.LOG = logger
	return logger, nil
}

// ProvideFiberApp 提供 Fiber App
func ProvideFiberApp(logger *zap.Logger) *fiber.App {
	return fiber.New(model.RunCONFIG.FiberConfig)
}
