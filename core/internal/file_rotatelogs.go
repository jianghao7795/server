package internal

import (
	"os"

	global "server/model"

	"go.uber.org/zap/zapcore"
)

type fileRotatelogs struct{}

var FileRotatelogs = new(fileRotatelogs)

// GetWriteSyncer 获取 zapcore.WriteSyncer
// Author wuhao
func (r *fileRotatelogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := NewCutter(global.CONFIG.Zap.Director, level, WithCutterFormat("2006-01-02"))
	if global.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
