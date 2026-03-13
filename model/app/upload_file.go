/*
 * @Author: jianghao
 * @Date: 2022-10-17 15:44:12
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 15:52:42
 */

package app

import global "server/model"

type FileUploadAndDownload struct {
	global.MODEL
	Name string `json:"name" form:"name" gorm:"comment:文件名"` // 文件名
	Url  string `json:"url" form:"url" gorm:"comment:文件地址"`  // 文件地址
	Tag  string `json:"tag" form:"tag" gorm:"comment:文件标签"`  // 文件标签
	Key  string `json:"key" form:"key" gorm:"comment:编号"`    // 编号
}

func (FileUploadAndDownload) TableName() string {
	return "exa_file_upload_and_downloads"
}
