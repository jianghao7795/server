package core

import (
	"os"
	"server/core/internal"
	global "server/model"
	"server/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap 获取 zap.Logger
// Author wuhao
func zapInit() (logger *zap.Logger, err error) {
	ok, err := utils.PathExists(global.CONFIG.Zap.Director)
	if err != nil {
		return
	}
	if !ok { // 判断是否有Director文件夹
		// log.Printf("create %v directory\n", global.CONFIG.Zap.Director)
		err = os.Mkdir(global.CONFIG.Zap.Director, os.ModePerm)
		if err != nil {
			return
		}
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return
}
