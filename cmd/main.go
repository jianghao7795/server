package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server-fiber/core"
	_ "server-fiber/docs" // 引入生成的文档
	global "server-fiber/model"
	"server-fiber/service/system"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
func main() {
	// 使用 Wire 初始化应用
	routerApp, cleanup, err := core.InitializeApp()
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}
	defer cleanup()

	app := routerApp.App

	// 设置全局日志
	zap.ReplaceGlobals(global.LOG)

	// 启动服务器
	go runServer(app)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		global.LOG.Error("Server Shutdown: ", zap.Error(err))
		os.Exit(1)
	}
	global.LOG.Info("Server exiting")
}

func runServer(app *fiber.App) {
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	global.LOG.Info("server run success on ", zap.String("address", address))
	log.Println(`Welcome to Fiber API`)

	if global.DB != nil {
		system.LoadAll() // 加载所有的 拉黑的jwt数据 避免盗用jwt
		// 程序结束前关闭数据库链接
		db, _ := global.DB.DB()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				global.LOG.Error("数据库关闭失败: " + err.Error())
			}
		}(db)
	}

	err := app.Listen(address)
	if err != nil {
		global.LOG.Error("Server run failed: " + err.Error())
		os.Exit(1)
	}
}
