package system

import (
	"errors"
	global "server/model"
	"server/model/system"
	"server/model/system/request"

	"gorm.io/gorm"
)

//
//@function: DeleteSysDictionary
//@description: 创建字典数据
//@param: sysDictionary model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) CreateSysDictionary(sysDictionary system.SysDictionary) (err error) {
	if (!errors.Is(global.DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound)) {
		return errors.New("存在相同的type，不允许创建")
	}
	err = global.DB.Create(&sysDictionary).Error
	return err
}

//
//@function: DeleteSysDictionary
//@description: 删除字典数据
//@param: sysDictionary model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) DeleteSysDictionary(id uint) (err error) {
	var sysDictionary system.SysDictionary
	err = global.DB.Where("id = ?", id).Preload("SysDictionaryDetails").First(&sysDictionary).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("请不要搞事(乱删)")
	}
	if err != nil {
		return err
	}
	err = global.DB.Delete(&sysDictionary).Error
	if err != nil {
		return err
	}

	if sysDictionary.SysDictionaryDetails != nil {
		return global.DB.Where("sys_dictionary_id=?", sysDictionary.ID).Delete(sysDictionary.SysDictionaryDetails).Error
	}
	return
}

//
//@function: UpdateSysDictionary
//@description: 更新字典数据
//@param: sysDictionary *model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) UpdateSysDictionary(sysDictionary *system.SysDictionary) (err error) {
	var dict system.SysDictionary
	sysDictionaryMap := map[string]any{
		"Name":   sysDictionary.Name,
		"Type":   sysDictionary.Type,
		"Status": sysDictionary.Status,
		"Desc":   sysDictionary.Desc,
	}
	db := global.DB.Where("id = ?", sysDictionary.ID).First(&dict)
	if dict.Type != sysDictionary.Type {
		if !errors.Is(global.DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同的type，不允许创建")
		}
	}
	err = db.Updates(sysDictionaryMap).Error
	return err
}

//
//@function: GetSysDictionary
//@description: 根据id或者type获取字典单条数据
//@param: Type string, Id uint
//@return: err error, sysDictionary model.SysDictionary

//	func (dictionaryService *DictionaryService) GetSysDictionary(Type string, Id uint) (sysDictionary system.SysDictionary, err error) {
//		err = global.DB.Where("type = ? OR id = ? and status = ?", Type, Id, true).Preload("SysDictionaryDetails", "status = ?", true).First(&sysDictionary).Error
//		return
//	}
func (dictionaryService *DictionaryService) GetSysDictionary(fieldType string, fieldID string) (sysDictionary system.SysDictionary, err error) {
	err = global.DB.Where("id = ? OR type = ? and status = ?", fieldID, fieldType, true).Preload("SysDictionaryDetails", "status = ?", true).First(&sysDictionary).Error
	return
}

//
//@author: wuhao
//@function: GetSysDictionaryInfoList
//@description: 分页获取字典列表
//@param: info request.SysDictionarySearch
//@return: err error, list any, total int64

func (dictionaryService *DictionaryService) GetSysDictionaryInfoList(info request.SysDictionarySearch) (list any, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&system.SysDictionary{})
	var sysDictionarys []system.SysDictionary
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("`name` LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != "" {
		db = db.Where("`type` LIKE ?", "%"+info.Type+"%")
	}
	if info.Status != nil {
		db = db.Where("`status` = ?", info.Status)
	}
	if info.Desc != "" {
		db = db.Where("`desc` LIKE ?", "%"+info.Desc+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&sysDictionarys).Error
	return sysDictionarys, total, err
}
