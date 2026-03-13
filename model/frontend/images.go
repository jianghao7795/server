package frontend

import global "server/model"

type ExaFile struct {
	global.MODEL
	Name       string  `json:"name" gorm:"comment:文件名"`        // 文件名
	Url        string  `json:"url" gorm:"comment:文件地址"`        // 文件地址
	Tag        string  `json:"tag" gorm:"comment:文件标签"`        // 文件标签
	Key        string  `json:"key" gorm:"comment:编号"`          // 编号
	Height     int     `json:"height" gorm:"comment:高度"`       // 图片高度
	Width      int     `json:"width" gorm:"comment:宽度"`        // 图片宽度
	Proportion float64 `json:"proportion" gorm:"comment:长宽比例"` // 图片长宽比例
}

func (ExaFile) TableName() string {
	return "exa_file_upload_and_downloads"
}
