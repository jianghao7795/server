package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server-fiber/init_load"
	global "server-fiber/model"
	"server-fiber/router"
	"server-fiber/service/system"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func initConfig() {
	var err error
	global.VIP, err = viperInit() // 初始化Viper 配置
	if err != nil {
		// global.LOG.Error("配置错误：", zap.Error(err))
		log.Println("配置错误： ", err.Error())
		os.Exit(1)
	}
	global.LOG, err = zapInit() // 初始化zap日志库
	if err != nil {
		log.Println("日志初始化错误： ", err.Error())
		// 日志初始化错误 退出程序
		os.Exit(1)
	}
	// global.Logger = InitLogger()   // 初始化 log 让log标准输出
	zap.ReplaceGlobals(global.LOG) // 配置部署到全局

	db, err := init_load.Gorm() // gorm连接数据库
	if err == nil {
		global.DB = db
		global.LOG.Info("Database connection success", zap.String("port", global.CONFIG.Mysql.Port))
	} else {
		global.LOG.Error("The database connection failed: " + err.Error())
		os.Exit(1)
	}
	init_load.Timer() // 定时 执行任务
	// err = utilsInit.TransInit("zh")
	// if err != nil {
	// 	global.LOG.Error("翻译错误：" + err.Error())
	// }

	if global.CONFIG.System.UseMultipoint || global.CONFIG.System.UseRedis {
		err = init_load.Redis()
		if err != nil {
			global.LOG.Error("Redis init failed: " + err.Error())
			os.Exit(1)
		}
	}

}

func RunServerElectron(app *fiber.App) {
	// 使用全局路由实例（向后兼容）
	appRouter := router.AppRouterInstance
	systemRouter := router.SystemRouterInstance
	exampleRouter := router.ExampleRouterInstance
	frontendRouter := router.FrontendRouterInstance
	mobileRouter := router.MobileRouterInstance
	routerApp := init_load.Routers(app, appRouter, systemRouter, exampleRouter, frontendRouter, mobileRouter)
	router := routerApp
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	global.LOG.Info("server run success on ", zap.String("address", address))
	log.Println(`Welcome to Fiber API`)
	err := router.Listen(address)
	if err != nil {
		global.LOG.Error("Server run failed: " + err.Error())
		os.Exit(1)
	}

	if global.DB != nil {
		system.LoadAll() // 加载所有的 拉黑的jwt数据 避免盗用jwt
		// initialize.RegisterTables(global.DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.DB.DB()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				global.LOG.Error("数据库关闭失败: " + err.Error())
			}
		}(db)
	}
}

func RunServer() {
	initConfig()
	app := fiber.New(global.RunCONFIG.FiberConfig)
	go RunServerElectron(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		global.LOG.Error("Server Shutdown: ", zap.Error(err))
		os.Exit(1)
	}
	global.LOG.Info("Server exiting")
}
