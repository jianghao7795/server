package system

import (
	global "server/model"
	"server/model/system"
	"server/model/system/request"
)

//
//@function: CreateSysDictionaryDetail
//@description: 创建字典详情数据
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) CreateSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	err = global.DB.Create(&sysDictionaryDetail).Error
	return err
}

//
//@function: DeleteSysDictionaryDetail
//@description: 删除字典详情数据
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) DeleteSysDictionaryDetail(id uint) (err error) {
	err = global.DB.Delete(&system.SysDictionaryDetail{}, id).Error
	return err
}

//
//@function: UpdateSysDictionaryDetail
//@description: 更新字典详情数据
//@param: sysDictionaryDetail *model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) UpdateSysDictionaryDetail(sysDictionaryDetail *system.SysDictionaryDetail) (err error) {
	err = global.DB.Model(&system.SysDictionaryDetail{}).Where("id = ?", sysDictionaryDetail.ID).Save(sysDictionaryDetail).Error
	return err
}

//
//@function: GetSysDictionaryDetail
//@description: 根据id获取字典详情单条数据
//@param: id uint
//@return: err error, sysDictionaryDetail model.SysDictionaryDetail

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetail(id uint) (sysDictionaryDetail system.SysDictionaryDetail, err error) {
	err = global.DB.Where("id = ?", id).First(&sysDictionaryDetail).Error
	return
}

//
//@function: GetSysDictionaryDetailInfoList
//@description: 分页获取字典详情列表
//@param: info request.SysDictionaryDetailSearch
//@return: err error, list any, total int64

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetailInfoList(info request.SysDictionaryDetailSearch) (list []system.SysDictionaryDetail, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&system.SysDictionaryDetail{})
	// utils.MergeQuery(db, info, "label", "value", "status", "sysDictionaryID")
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Label != "" {
		db = db.Where("label LIKE ?", "%"+info.Label+"%")
	}
	if info.Value != 0 {
		db = db.Where("value = ?", info.Value)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysDictionaryID != 0 {
		db = db.Where("sys_dictionary_id = ?", info.SysDictionaryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}
