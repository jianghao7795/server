package utils

import (
	"os"

	global "server/model"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

//@function: GetWriteSyncer
//@description: zap logger中加入file-rotatelogs
//@return: zapcore.WriteSyncer, error

func GetWriteSyncer(file string) zapcore.WriteSyncer {
	s, err := os.Stat(file)
	if err != nil {
		return zapcore.AddSync(nil)
	}
	if s.IsDir() {
		return zapcore.AddSync(nil)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    1,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}

	if global.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
