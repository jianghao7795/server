package example

import (
	v1 "server/api/v1/example"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type FileUploadAndDownloadRouter struct{}

func (e *FileUploadAndDownloadRouter) InitFileUploadAndDownloadRouter(Router fiber.Router) {
	fileUploadAndDownloadRouter := Router.Group("fileUploadAndDownload")
	exaFileUploadAndDownloadApi := new(v1.FileUploadAndDownloadApi)

	fileUploadAndDownloadRouter.Post("upload", middleware.OperationRecord, exaFileUploadAndDownloadApi.UploadFile)                                 // 上传文件
	fileUploadAndDownloadRouter.Get("getFileList", exaFileUploadAndDownloadApi.GetFileList)                                                        // 获取上传文件列表
	fileUploadAndDownloadRouter.Delete("deleteFile/:id", middleware.OperationRecord, exaFileUploadAndDownloadApi.DeleteFile)                       // 删除指定文件
	fileUploadAndDownloadRouter.Put("editFileName", middleware.OperationRecord, exaFileUploadAndDownloadApi.EditFileName)                          // 编辑文件名或者备注
	fileUploadAndDownloadRouter.Post("breakpointContinue", middleware.OperationRecord, exaFileUploadAndDownloadApi.BreakpointContinue)             // 断点续传
	fileUploadAndDownloadRouter.Get("findFile", exaFileUploadAndDownloadApi.FindFile)                                                              // 查询当前文件成功的切片
	fileUploadAndDownloadRouter.Post("breakpointContinueFinish", middleware.OperationRecord, exaFileUploadAndDownloadApi.BreakpointContinueFinish) // 切片传输完成
	fileUploadAndDownloadRouter.Delete("removeChunk", middleware.OperationRecord, exaFileUploadAndDownloadApi.RemoveChunk)                         // 删除切片
	fileUploadAndDownloadRouter.Get("getFileBreakpoint", exaFileUploadAndDownloadApi.FindFileBreakpoint)                                           // 查询当前文件所有 断点上车文件
	fileUploadAndDownloadRouter.Delete("deleteFileBreakpoint/:id", middleware.OperationRecord, exaFileUploadAndDownloadApi.DeleteFileBreakpoint)   // 删除断点续传文件
}
