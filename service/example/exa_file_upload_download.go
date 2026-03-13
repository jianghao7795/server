package example

import (
	"errors"
	"mime/multipart"
	"server/model/common/request"
	"server/model/example"
	"server/utils/upload"
	"strings"

	global "server/model"

	fileDimensionReq "server/model/example/request"
)

//@author: wuhao
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) Upload(file example.ExaFileUploadAndDownload) error {
	return global.DB.Create(&file).Error
}

//@author: wuhao
//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: error, model.ExaFileUploadAndDownload

func (e *FileUploadAndDownloadService) FindFile(id uint) (example.ExaFileUploadAndDownload, error) {
	var file example.ExaFileUploadAndDownload
	err := global.DB.Where("id = ?", id).First(&file).Error
	return file, err
}

//@author: wuhao
//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownloadService) DeleteFile(id uint) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	fileFromDb, err = e.FindFile(id)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Url); err != nil {
		return errors.New("文件删除失败 " + err.Error())
	}
	err = global.DB.Delete(&fileFromDb, id).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file *example.ExaFileUploadAndDownload) (err error) {
	return global.DB.Where("id = ?", file.ID).Update("name", file.Name).Error
}

//@author: wuhao
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list any, total int64

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info *request.PageInfo) (list []example.ExaFileUploadAndDownload, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := global.DB.Model(&example.ExaFileUploadAndDownload{})
	var fileLists []example.ExaFileUploadAndDownload
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if len(info.IsCropper) == 1 {
		db = db.Where("is_cropper = ?", info.IsCropper)
	}

	if info.Proportion != "" {
		db = db.Where("proportion >= ?", info.Proportion)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&fileLists).Error
	return fileLists, total, err
}

//@author: wuhao
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: err error, file model.ExaFileUploadAndDownload

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string, fileDimension *fileDimensionReq.FileDimension, isCropper int) (file example.ExaFileUploadAndDownload, err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return example.ExaFileUploadAndDownload{}, uploadErr
	}
	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		f := example.ExaFileUploadAndDownload{
			Url:        filePath,
			Name:       header.Filename,
			Tag:        s[len(s)-1],
			Key:        key,
			Width:      fileDimension.Width,
			Height:     fileDimension.Height,
			Proportion: fileDimension.Proportion,
			IsCropper:  isCropper,
		}
		return f, e.Upload(f)
	}
	return
}
