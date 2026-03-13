/*
 * @Author: jianghao
 * @Date: 2022-10-17 11:13:10
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 17:15:16
 */

package app

import (
	v1 "server/api/v1/app"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type BaseMessageRouter struct{}

// InitArticleRouter 初始化 base_message 路由信息
func (r *BaseMessageRouter) InitBaseMessageRouter(c fiber.Router) {
	baseMessageRouter := c.Group("base_message")
	var baseMessageApi = new(v1.BaseMessageApi)
	var uploadFileApi = new(v1.FileUploadAndDownloadApi)

	baseMessageRouter.Post("createBaseMessage", middleware.OperationRecord, baseMessageApi.CreateBaseMessage)
	baseMessageRouter.Put("updateBaseMessage", middleware.OperationRecord, baseMessageApi.UpdateBaseMessage)
	baseMessageRouter.Post("upload_file", middleware.OperationRecord, uploadFileApi.UploadFile)

	baseMessageRouter.Get("getBaseMessage/:id", baseMessageApi.FindBaseMessage)

}
