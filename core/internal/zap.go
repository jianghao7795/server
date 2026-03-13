package internal

import (
	"time"

	global "server/model"

	gozap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap = new(zap)

type zap struct{}

// GetEncoder 获取 zapcore.Encoder
// Author wuhao
func (z *zap) GetEncoder() zapcore.Encoder {
	if global.CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

// GetEncoderConfig 获取zapcore.EncoderConfig
// Author wuhao
func (z *zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    global.CONFIG.Zap.ZapEncodeLevel(),
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // FullCallerEncoder
	}
}

// GetEncoderCore 获取Encoder的 zapcore.Core
// Author wuhao
func (z *zap) GetEncoderCore(l zapcore.Level, level gozap.LevelEnablerFunc) zapcore.Core {
	writer := FileRotatelogs.GetWriteSyncer(l.String()) // 日志分割
	return zapcore.NewCore(z.GetEncoder(), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
// Author wuhao
func (z *zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(global.CONFIG.Zap.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core
// Author wuhao
func (z *zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := global.CONFIG.Zap.TransportLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	return cores
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
// Author wuhao
func (z *zap) GetLevelPriority(level zapcore.Level) gozap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == gozap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == gozap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == gozap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == gozap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == gozap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == gozap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == gozap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == gozap.DebugLevel
		}
	}
}
