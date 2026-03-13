package utils

import (
	global "server/model"

	"go.uber.org/zap"
)

func HandleError(err error, message string) bool {
	if err != nil {
		global.LOG.Error(message, zap.Error(err))
		return false
	}
	return true
}
