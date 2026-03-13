package example

import (
	"errors"
	"os"
	global "server/model"
	"server/model/example"

	"gorm.io/gorm"
)

type FileUploadAndDownloadService struct{}

//@author: wuhao
//@function: FindOrCreateFile
//@description: 上传文件时检测当前文件属性，如果没有文件则创建，有则返回文件的当前切片
//@param: fileMd5 string, fileName string, chunkTotal int
//@return: err error, file model.ExaFile

func (e *FileUploadAndDownloadService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (file example.ExaFile, err error) {
	var cfile example.ExaFile
	cfile.FileMd5 = fileMd5
	cfile.FileName = fileName
	cfile.ChunkTotal = chunkTotal

	if errors.Is(global.DB.Where("file_md5 = ? AND is_finish = ?", fileMd5, true).First(&file).Error, gorm.ErrRecordNotFound) {
		err = global.DB.Where("file_md5 = ? AND file_name = ?", fileMd5, fileName).Preload("ExaFileChunk").FirstOrCreate(&file, cfile).Error
		return file, err
	}
	cfile.IsFinish = true
	cfile.FilePath = file.FilePath
	err = global.DB.Create(&cfile).Error
	return cfile, err
}

//@author: wuhao
//@function: CreateFileChunk
//@description: 创建文件切片记录
//@param: id uint, fileChunkPath string, fileChunkNumber int
//@return: error

func (e *FileUploadAndDownloadService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	var chunk example.ExaFileChunk
	chunk.FileChunkPath = fileChunkPath
	chunk.ExaFileID = id
	chunk.FileChunkNumber = fileChunkNumber
	err := global.DB.Create(&chunk).Error
	return err
}

//@author: wuhao
//@function: DeleteFileChunk
//@description: 删除文件切片记录
//@param: fileMd5 string, fileName string, filePath string
//@return: error

func (e *FileUploadAndDownloadService) DeleteFileChunk(fileMd5 string, fileName string, filePath string) error {
	var chunks []example.ExaFileChunk
	var file example.ExaFile
	err := global.DB.Where("file_md5 = ? ", fileMd5).First(&file).Update("is_finish", true).Update("file_path", filePath).Error
	if err != nil {
		return err
	}
	err = global.DB.Where("exa_file_id = ?", file.ID).Delete(&chunks).Unscoped().Error
	return err
}

// @author: wuhao
// @function: DeleteFileBreakpoint
// @description: 删除断点上传文件
// @param: id uint
// @return: error
func (e *FileUploadAndDownloadService) DeleteFileBreakpoint(id int) (err error) {
	var file example.ExaFile
	db := global.DB.Where("id = ?", id).First(&file)
	if file.IsFinish {
		err := os.Remove(file.FilePath)
		if err != nil {
			return err
		}
	}

	err = db.Delete(&file).Error
	return
}

// @author: wuhao
// @function: FindFileBreakpoint
// @description: 断点上传文件列表
// @param: info request.Empty
// @return: []example.ExaFile, error
func (e *FileUploadAndDownloadService) FindFileBreakpoint(page int, pageSize int) ([]example.ExaFile, int64, error) {
	var files []example.ExaFile
	offset := pageSize * (page - 1)
	var total int64
	db := global.DB.Model(&example.ExaFile{})
	db = db.Count(&total)
	err := db.Offset(offset).Limit(pageSize).Order("id desc").Find(&files).Error
	return files, total, err
}
