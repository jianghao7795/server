package example

import (
	"os"
	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/example"
	"server/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// /excel/importExcel 接口，与upload接口作用类似，只是把文件存到resource/excel目录下，用于导入Excel时存放Excel文件(ExcelImport.xlsx)
// /excel/loadExcel接口，用于读取resource/excel目录下的文件((ExcelImport.xlsx)并加载为[]model.SysBaseMenu类型的示例数据
// /excel/exportExcel 接口，用于读取前端传来的tableData，生成Excel文件并返回
// /excel/downloadTemplate 接口，用于下载resource/excel目录下的 ExcelTemplate.xlsx 文件，作为导入的模板

// ExportExcel 导出Excel
// @Tags excel
// @Summary 导出Excel
// @Description 根据提供的数据导出Excel文件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/octet-stream
// @Param data body example.ExcelInfo true "导出Excel文件信息"
// @Success 200 {file} file "Excel文件"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /excel/exportExcel [post]
func (e *ExcelApi) ExportExcel(c *fiber.Ctx) (err error) {
	excelInfo := example.ExcelInfo{
		FileName: "",
		InfoList: nil,
	}
	if err := c.BodyParser(&excelInfo); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	// log.Println("excelInfo: ", excelInfo.InfoList)

	if excelInfo.FileName == "" {
		return response.FailWithMessage("请传参数filename", c)
	}
	if strings.Contains(excelInfo.FileName, "..") {
		return response.FailWithMessage("包含非法字符", c)
	}
	// err := excelService.GetMenuData(&excelInfo.InfoList)
	// if err != nil {
	// 	global.LOG.Error("获取失败!", zap.Error(err))
	// 	return response.FailWithMessage("获取失败", c)
	// }
	filePath := global.CONFIG.Excel.Dir + excelInfo.FileName
	err = excelService.ParseInfoList2Excel(excelInfo.InfoList, filePath)
	if err != nil {
		global.LOG.Error("转换Excel失败!", zap.Error(err))
		return response.FailWithMessage("转换Excel失败", c)
	}
	c.Set("success", "true")
	return c.SendFile(filePath)
}

// ImportExcel 导入Excel文件
// @Tags excel
// @Summary 导入Excel文件
// @Description 上传并导入Excel文件到系统
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {object} response.Response{msg=string} "导入Excel文件成功"
// @Failure 400 {object} response.Response "文件格式错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /excel/importExcel [post]
func (e *ExcelApi) ImportExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		global.LOG.Error("接收文件失败!", zap.Error(err))
		return response.FailWithMessage("接收文件失败", c)
	}
	filepath_time := time.Now().Format("2006/01/02")
	filenameMd5 := utils.MD5V([]byte(file.Filename)) + "_" + time.Now().Format("20060102150405") + "." + strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
	fileTypeName := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
	importExcel := example.FielUploadImport{
		FileName:    file.Filename,
		FileNameMd5: filenameMd5,
		State:       1,
		FileSize:    file.Size,
		FilePath:    global.CONFIG.Excel.Dir + filepath_time + "/" + filenameMd5,
		FileType:    fileTypeName,
	}
	err = excelService.ImportExcel(&importExcel)
	if err != nil {
		global.LOG.Error("导入失败!", zap.Error(err))
		return response.FailWithMessage("导入失败", c)
	}
	mkdirErr := os.MkdirAll(global.CONFIG.Excel.Dir+filepath_time, os.ModePerm)
	if mkdirErr != nil {
		global.LOG.Error("创建目录失败：", zap.Any("err", mkdirErr.Error()))
	}
	err = c.SaveFile(file, global.CONFIG.Excel.Dir+filepath_time+"/"+filenameMd5)
	if err != nil {
		global.LOG.Error("保存文件失败：", zap.Any("err", err.Error()))
	}
	return response.OkWithMessage("导入成功", c)
}

// @Tags excel
// @Summary 加载Excel数据
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{data=object,msg=string} "加载Excel数据,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /excel/loadExcel [get]
func (e *ExcelApi) LoadExcel(c *fiber.Ctx) error {
	menus, err := excelService.ParseExcel2InfoList()
	if err != nil {
		global.LOG.Error("加载数据失败!", zap.Error(err))
		return response.FailWithMessage("加载数据失败", c)
	}
	return response.OkWithDetailed(response.PageResult{
		List:     menus,
		Total:    int64(len(menus)),
		Page:     1,
		PageSize: 999,
	}, "加载数据成功", c)
}

// @Tags excel
// @Summary 下载模板
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param fileName query string true "模板名称"
// @Success 200 {object} response.Response{msg=string} "下载模板"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /excel/downloadTemplate [get]
func (e *ExcelApi) DownloadTemplate(c *fiber.Ctx) error {
	fileName := c.Query("fileName")
	filePath := global.CONFIG.Excel.Dir + fileName

	fi, err := os.Stat(filePath)
	if err != nil {
		global.LOG.Error("文件不存在!", zap.Error(err))
		return response.FailWithMessage("文件不存在", c)
	}
	if fi.IsDir() {
		global.LOG.Error("不支持下载文件夹!", zap.Error(err))
		return response.FailWithMessage("不支持下载文件夹", c)
	}
	c.Set("success", "true") // 增加返回头 信息
	return c.Download(filePath)
}

func (e *ExcelApi) GetFileList(c *fiber.Ctx) error {
	var pageInfo request.PageInfo
	_ = c.QueryParser(&pageInfo)
	list, total, err := excelService.GetFileList(&pageInfo)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (e *ExcelApi) DeleteFile(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := excelService.DeleteFile(int64(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}
