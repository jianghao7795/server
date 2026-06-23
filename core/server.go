package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	global "server/model"
	"server/service/system"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func RunServer() {
	routerApp, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatalf("initialize Fiber app failed: %v", err)
	}
	defer cleanup()
	defer closeDatabase()

	app := routerApp.App
	zap.ReplaceGlobals(global.LOG)

	if global.DB != nil {
		system.LoadAll()
	}

	go listen(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		global.LOG.Error("shutdown Fiber app failed", zap.Error(err))
		os.Exit(1)
	}
}

func listen(app *fiber.App) {
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	log.Println("Welcome to Fiber API")

	if err := app.Listen(address); err != nil {
		global.LOG.Error("listen Fiber app failed", zap.Error(err))
		os.Exit(1)
	}
}

func closeDatabase() {
	if global.DB == nil {
		return
	}

	db, err := global.DB.DB()
	if err != nil {
		global.LOG.Error("get database handle failed", zap.Error(err))
		return
	}

	if err := db.Close(); err != nil {
		global.LOG.Error("close database failed", zap.Error(err))
	}
}
