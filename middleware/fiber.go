package middleware

import (
	"runtime/debug"

	global "server/model"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
)

func AppMiddlewares() []any {
	return []any{
		Recovery(),
		RequestID(),
		AccessLogger(),
		ResponseCompressor(),
		DefaultLimit,
		CorsByRules,
	}
}

func Recovery() fiber.Handler {
	return fiberrecover.New(fiberrecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, err any) {
			if global.LOG == nil {
				return
			}
			global.LOG.Error("panic recovered",
				zap.Any("panic", err),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
				zap.ByteString("stack", debug.Stack()),
			)
		},
	})
}

func RequestID() fiber.Handler {
	return requestid.New()
}

func AccessLogger() fiber.Handler {
	return logger.New(logger.Config{
		Done: global.Done,
	})
}

func ResponseCompressor() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
