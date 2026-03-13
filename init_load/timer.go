package init_load

import (
	"server/config"
	global "server/model"
	"server/utils"

	"go.uber.org/zap"
)

func Timer() {
	if global.CONFIG.Timer.Start {
		for i := range global.CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				_, err := global.Timer.AddTaskByFunc("UpdateGithub", global.CONFIG.Timer.Spec, func() {
					err := utils.UpdateTable(global.DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						global.LOG.Error("新增Github记录错误：", zap.Error(err))
					}
				})
				if err != nil {
					return
				}
			}(global.CONFIG.Timer.Detail[i])
		}
	}
}
