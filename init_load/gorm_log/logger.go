package gorm_log

import (
	"fmt"

	global "server/model"

	"gorm.io/gorm/logger"
)

type Writer struct {
	logger.Writer
}

// NewWriter writer 构造函数

func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

// Printf 格式化打印日志

func (w *Writer) Printf(message string, data ...any) {
	var logZap bool
	switch global.CONFIG.System.DbType {
	case "mysql":
		logZap = global.CONFIG.Mysql.LogZap
	case "pgsql":
		logZap = global.CONFIG.Pgsql.LogZap
	}
	if logZap {
		global.LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
