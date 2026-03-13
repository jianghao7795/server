package app

import (
	global "server/model"
	"server/model/common/response"

	"github.com/gofiber/fiber/v2"
)

// StartTasking 启动任务
// @Tags Tasking
// @Summary 启动任务
// @Description 启动指定的定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param task query string true "任务名称"
// @Success 200 {object} response.Response{msg=string} "启动任务成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tasking/start [get]
func (*TaskNameApi) StartTasking(c *fiber.Ctx) error {
	tasking := c.Query("task")
	if tasking == "" {
		global.LOG.Error("请传入任务名!")
		return response.FailWithMessage("请传入任务名", c)
	}
	_, status := global.Timer.FindCron(tasking)

	if !status {
		global.LOG.Error("开启失败!")
		return response.FailWithMessage("开启失败，没有这个任务", c)
	} else {
		global.Timer.StartTask(tasking)
		return response.OkWithMessage("开启成功", c)
	}
}

// StopTasking 停止任务
// @Tags Tasking
// @Summary 停止任务
// @Description 停止指定的定时任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param task query string true "任务名称"
// @Success 200 {object} response.Response{msg=string} "停止任务成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tasking/stop [get]
func (*TaskNameApi) StopTasking(c *fiber.Ctx) error {
	tasking := c.Query("task")
	if tasking == "" {
		global.LOG.Error("请传入任务名!")
		return response.FailWithMessage("请传入任务名", c)
	}
	_, status := global.Timer.FindCron(tasking)

	if !status {
		global.LOG.Error("关闭失败!")
		return response.FailWithMessage("关闭失败，没有这个任务", c)
	} else {
		global.Timer.StopTask(tasking)
		return response.OkWithMessage("关闭成功", c)
	}
}
