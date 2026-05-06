package example

import (
	ioutil "io"
	"mime/multipart"
	"strconv"

	"server/model/common/response"
	"server/model/example"
	"server/model/example/request"
	"server/utils"

	exampleRes "server/model/example/response"

	"github.com/gofiber/fiber/v3"
)

// BreakpointContinue @Tags ExaFileUploadAndDownload
// @Summary 断点续传到服务器
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "an example for breakpoint resume, 断点续传示例"
// @Success 200 {object} response.Response{msg=string} "断点续传到服务器"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/breakpointContinue [post]
func (u *FileUploadAndDownloadApi) BreakpointContinue(c fiber.Ctx) error {
	var breakpoint example.ExaFileData
	err := c.Bind().Body(&breakpoint)
	if err != nil {
		return response.FailWithMessage("获取数据失败", 3, err, c)
	}
	breakpoint.FileHeader, err = c.FormFile("file")
	if err != nil {
		return response.FailWithMessage("接收文件失败: "+err.Error(), 3, err, c)
	}
	f, err := breakpoint.FileHeader.Open()
	if err != nil {
		return response.FailWithMessage("文件读取失败"+err.Error(), 3, err, c)
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)

	cen, err := ioutil.ReadAll(f)
	if err != nil {
		return response.FailWithMessage("文件分段读取失败", 3, err, c)
	}
	if !utils.CheckMd5(cen, breakpoint.ChunkMd5) {
		return response.FailWithMessage("检查md5失败", 3, err, c)
	}
	file, err := fileUploadAndDownloadService.FindOrCreateFile(breakpoint.FileMd5, breakpoint.FileName, breakpoint.ChunkTotal)
	if err != nil {
		return response.FailWithMessage("查找或创建记录失败", 3, err, c)
	}
	paths, err := utils.BreakPointContinue(cen, breakpoint.FileName, breakpoint.ChunkNumber, breakpoint.ChunkTotal, breakpoint.FileMd5)
	if err != nil {
		return response.FailWithMessage("断点续传失败: "+err.Error(), 3, err, c)
	}

	if err = fileUploadAndDownloadService.CreateFileChunk(file.ID, paths, breakpoint.ChunkNumber); err != nil {
		return response.FailWithMessage("创建文件记录失败: "+err.Error(), 3, err, c)
	}
	return response.OkWithMessage("切片创建成功", c)
}

// @Tags ExaFileUploadAndDownload
// @Summary 查找文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "Find the file, 查找文件"
// @Success 200 {object} response.Response{msg=string} "查找文件,返回包括文件详情"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/findFile [get]
func (u *FileUploadAndDownloadApi) FindFile(c fiber.Ctx) error {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	chunkTotal, _ := strconv.Atoi(c.Query("chunkTotal", "0"))
	if chunkTotal == 0 {
		return response.FailWithMessage("获取文件大小失败", 3, nil, c)
	}
	file, err := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		return response.FailWithMessage("查找失败: "+err.Error(), 3, err, c)
	} else {
		return response.OkWithDetailed(exampleRes.FileResponse{File: file}, "查找成功", c)
	}
}

// @Tags ExaFileUploadAndDownload
// @Summary 创建文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件完成"
// @Success 200 {object} response.Response{msg=string} "创建文件,返回包括文件路径"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/findFile [post]
func (b *FileUploadAndDownloadApi) BreakpointContinueFinish(c fiber.Ctx) error {
	var file request.BreakPoint
	// filepath := c.Body("filePath")
	err := c.Bind().Body(&file)
	if err != nil {
		return response.FailWithMessage("获取文件信息错误: "+err.Error(), 3, err, c)
	}
	// log.Println("filename: ", file.FileName, " fileMd5: ", file.FileMd5)
	filePath, err := utils.MakeFile(file.FileName, file.FileMd5)
	if err != nil {
		return response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: filePath}, "文件创建失败", 3, err, c)
	}
	// err = fileUploadAndDownloadService.DeleteFileChunk(file.FileMd5, file.FileName, filePath)
	// if err != nil {
	// 	global.LOG.Error("删除切片失败", zap.Error(err))
	// 	return response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: filePath}, "删除切片失败", nil, c)
	// }

	return response.OkWithDetailed(exampleRes.FilePathResponse{FilePath: filePath}, "文件创建成功", c)
}

// @Tags ExaFileUploadAndDownload
// @Summary 删除切片
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "删除缓存切片"
// @Success 200 {object} response.Response{msg=string} "删除切片"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/removeChunk [delete]
func (u *FileUploadAndDownloadApi) RemoveChunk(c fiber.Ctx) error {
	var file example.ExaFile
	err := c.Bind().Body(&file)
	if err != nil {
		return response.FailWithMessage("缓存切片删除失败", 3, err, c)
	}
	err = utils.RemoveChunk(file.FileMd5)
	if err != nil {
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "缓存切片文件删除失败", 3, err, c)
	}
	err = fileUploadAndDownloadService.DeleteFileChunk(file.FileMd5, file.FileName, file.FilePath)
	if err != nil {
		return response.FailWithMessage(err.Error(), 3, err, c)
	} else {
		return response.OkWithMessage("表缓存切片删除成功", c)
	}
}

// @Tags ExaFileUploadAndDownload
// @Summary 查找缓存文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "Find the file, 查找缓存文件列表"
// @Success 200 {object} response.Response{data=object,msg=string} "查找缓存文件列表"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/findFileBreakpoint [get]
func (u *FileUploadAndDownloadApi) FindFileBreakpoint(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	file, total, err := fileUploadAndDownloadService.FindFileBreakpoint(page, pageSize)
	if err != nil {
		return response.FailWithMessage("获取失败: "+err.Error(), 3, err, c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     file,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, "查找成功", c)
	}
}

// @Tags ExaFileUploadAndDownload
// @Summary 删除文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "Find the file, 删除缓存文件"
// @Success 200 {object} response.Response{msg=string} "删除缓存文件"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /fileUploadAndDownload/deleteFileBreakpoint [delete]
func (u *FileUploadAndDownloadApi) DeleteFileBreakpoint(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取id失败: "+err.Error(), 3, err, c)
	}
	err = fileUploadAndDownloadService.DeleteFileBreakpoint(id)
	if err != nil {
		return response.FailWithMessage("删除失败: "+err.Error(), 3, err, c)
	} else {
		return response.OkWithMessage("删除失败", c)
	}
}
