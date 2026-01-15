package init_load

import (
	global "server-fiber/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RedisInitialized Redis 初始化完成标记
type RedisInitialized struct{}

// TimerInitialized 定时器初始化完成标记
type TimerInitialized struct{}

// RouterApp 已配置路由的 App
type RouterApp struct {
	App *fiber.App
}

// ProvideGorm 提供 GORM 数据库连接
func ProvideGorm(logger *zap.Logger) (*gorm.DB, error) {
	db, err := Gorm()
	if err != nil {
		return nil, err
	}
	global.DB = db
	return db, nil
}

// ProvideRedis 初始化 Redis（如果配置了）
func ProvideRedis(logger *zap.Logger) (RedisInitialized, error) {
	if global.CONFIG.System.UseMultipoint || global.CONFIG.System.UseRedis {
		if err := Redis(); err != nil {
			return RedisInitialized{}, err
		}
		return RedisInitialized{}, nil
	}
	return RedisInitialized{}, nil
}

// ProvideTimer 初始化定时器任务
func ProvideTimer(db *gorm.DB, logger *zap.Logger) TimerInitialized {
	Timer()
	return TimerInitialized{}
}

// ProvideRouter 初始化路由
func ProvideRouter(app *fiber.App, db *gorm.DB, logger *zap.Logger, _ RedisInitialized, _ TimerInitialized) *RouterApp {
	configuredApp := Routers(app)
	return &RouterApp{App: configuredApp}
}
