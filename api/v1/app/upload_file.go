/*
 * @Author: jianghao
 * @Date: 2022-10-17 15:45:55
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 17:04:47
 */

package app

import (
	global "server-fiber/model"
	responseUploadFile "server-fiber/model/app/response"
	"server-fiber/model/common/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UploadFile
// @Tags backend upload
// @Summary 上传文件示例
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件示例"
// @Success 200 {object} response.Response{data=responseUploadFile.ResponseUploadFile{file=app.FileUploadAndDownload},msg=string} "上传文件示例,返回包括文件详情"
// @Router /base_message/upload [post]
func (u *FileUploadAndDownloadApi) UploadFile(c *fiber.Ctx) error {
	// var file app.FileUploadAndDownload
	noSave := c.Query("noSave", "0")
	header, err := c.FormFile("file")
	if err != nil {
		global.LOG.Error("接收文件失败!", zap.Error(err))
		return response.FailWithMessage400("接收文件失败", c)
	}
	file, err := fileUploadService.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		global.LOG.Error("上传更新失败!", zap.Error(err))
		return response.FailWithMessage400("上传更新失败", c)
	}
	return response.OkWithDetailed(responseUploadFile.ResponseUploadFile{File: file}, "上传成功", c)
}
