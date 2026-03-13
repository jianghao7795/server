package example

import (
	"errors"
	"fmt"
	"os"
	global "server/model"
	"server/model/common/request"
	"server/model/example"
	"server/model/system"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

func (exa *ExcelService) ParseInfoList2Excel(infoList []example.FielUploadImport, filePath string) error {
	excel := excelize.NewFile()
	excel.SetSheetRow("router", "A1", &[]string{"ID", "FilePath", "FileSize", "FileType", "Filename", "FilenameMD5", "是否启用"})
	b := 0
	for i, fileMerge := range infoList {
		axis := fmt.Sprintf("A%d", i+b+2)
		excel.SetSheetRow("router", axis, &[]any{
			fileMerge.ID,
			fileMerge.FilePath,
			fileMerge.FileSize,
			fileMerge.FileType,
			fileMerge.FileName,
			fileMerge.FileNameMd5,
			fileMerge.State,
		})
		// log.Println("i: ", i)
		// if menu.Children != nil {
		// 	for _, m := range menu.Children {
		// 		b += 1
		// 		a := fmt.Sprintf("A%d", i+b+2)
		// 		excel.SetSheetRow("router", a, &[]any{
		// 			m.ID,
		// 			m.Name,
		// 			m.Path,
		// 			m.Hidden,
		// 			m.ParentId,
		// 			m.Sort,
		// 			m.Component,
		// 		})

		// 	}

		// }
	}
	// log.Println("filepath: ", filePath)
	err := excel.SaveAs(filePath)
	return err
}

func (exa *ExcelService) ParseExcel2InfoList() ([]system.SysBaseMenu, error) {
	skipHeader := true
	fixedHeader := []string{"ID", "路由Name", "路由Path", "是否隐藏", "父节点", "排序", "文件名称"}
	file, err := excelize.OpenFile(global.CONFIG.Excel.Dir + "ExcelImport.xlsx")
	if err != nil {
		return nil, err
	}
	menus := make([]system.SysBaseMenu, 0)
	rows, err := file.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if skipHeader {
			if exa.compareStrSlice(row, fixedHeader) {
				skipHeader = false
				continue
			} else {
				return nil, errors.New("Excel格式错误")
			}
		}
		if len(row) != len(fixedHeader) {
			continue
		}
		id, _ := strconv.Atoi(row[0])
		hidden, _ := strconv.ParseBool(row[3])
		sort, _ := strconv.Atoi(row[5])
		menu := system.SysBaseMenu{
			MODEL: global.MODEL{
				ID: uint(id),
			},
			Name:      row[1],
			Path:      row[2],
			Hidden:    hidden,
			ParentId:  row[4],
			Sort:      sort,
			Component: row[6],
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (exa *ExcelService) compareStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if (b == nil) != (a == nil) {
		return false
	}
	for key, value := range a {
		if value != b[key] {
			return false
		}
	}
	return true
}

func (exa *ExcelService) ImportExcel(data *example.FielUploadImport) error {
	return global.DB.Create(data).Error
}

func (exa *ExcelService) GetFileList(pageInfo *request.PageInfo) (lists []example.FielUploadImport, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	db := global.DB.Model(&example.FielUploadImport{})
	var fileLists []example.FielUploadImport
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&fileLists).Error
	if err != nil {
		return nil, 0, err
	}
	return fileLists, total, nil
}

func (exa *ExcelService) DeleteFile(id int64) error {
	var file example.FielUploadImport
	err := global.DB.Model(&example.FielUploadImport{}).Where("id = ?", id).First(&file).Error
	if err != nil {
		return err
	}
	err = global.DB.Delete(&file).Error
	if err != nil {
		return err
	}
	projectPath, _ := os.Getwd()
	err = os.Rename(projectPath+"/"+file.FilePath, projectPath+"/"+file.FilePath+".bak")
	if err != nil {
		return err
	}
	return nil
}

func (exa *ExcelService) GetMenuData(lists *[]system.SysBaseMenu) (err error) {
	db := global.DB.Model(&system.SysBaseMenu{})
	return db.Find(lists).Error
}
