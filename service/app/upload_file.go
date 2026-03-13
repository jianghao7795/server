/*
 * @Author: jianghao
 * @Date: 2022-10-17 15:48:45
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 15:51:35
 */

package app

import (
	"mime/multipart"
	global "server/model"
	"server/model/app"
	"server/utils/upload"
	"strings"
)

type FileUploadService struct{}

//@author: wuhao
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: err error, file model.ExaFileUploadAndDownload

func (e *FileUploadService) UploadFile(header *multipart.FileHeader, noSave string) (file app.FileUploadAndDownload, err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return app.FileUploadAndDownload{}, uploadErr
	}
	if noSave == "0" {
		last := strings.LastIndex(header.Filename, ".")
		s := header.Filename[last:]
		f := app.FileUploadAndDownload{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s,
			Key:  key,
		}
		return f, e.createdFileRecord(&f)
	}
	return
}

//@author: wuhao
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadService) createdFileRecord(file *app.FileUploadAndDownload) error {
	return global.DB.Create(file).Error
}
