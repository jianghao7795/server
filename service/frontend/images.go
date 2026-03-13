package frontend

import (
	global "server/model"
	"server/model/frontend"
)

type Images struct{}

func (s *Images) GetImagesList() (list []frontend.ExaFile, err error) {
	err = global.DB.Model(&frontend.ExaFile{}).Where("proportion > 1.6").Order("id desc").Find(&list).Error
	return
}
