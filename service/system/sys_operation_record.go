package system

import (
	global "server/model"
	"server/model/common/request"
	"server/model/system"
	systemReq "server/model/system/request"
)

//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord *system.SysOperationRecord) (err error) {
	err = global.DB.Create(sysOperationRecord).Error
	return err
}

//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids request.IdsReq
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecordByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecord(id uint) (err error) {
	err = global.DB.Delete(&system.SysOperationRecord{}, id).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: err error, sysOperationRecord model.SysOperationRecord

func (operationRecordService *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord system.SysOperationRecord, err error) {
	err = global.DB.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info systemReq.SysOperationRecordSearch
//@return: err error, list any, total int64

func (operationRecordService *OperationRecordService) GetSysOperationRecordInfoList(info *systemReq.SysOperationRecordSearch) (list []system.SysOperationRecord, total int64, err error) {
	// time.Sleep(3 * time.Second)
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&system.SysOperationRecord{})
	var sysOperationRecords []system.SysOperationRecord
	db = db.Where("type_port = ?", info.TypePort)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return sysOperationRecords, total, err
}

func (operationRecordService *OperationRecordService) GetSysOperationRecordInfoFrontendList(info systemReq.SysOperationRecordSearch) (list []system.SysOperationRecordFrontend, total int64, err error) {
	// time.Sleep(3 * time.Second)
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&system.SysOperationRecordFrontend{})

	var sysOperationRecords []system.SysOperationRecordFrontend
	db = db.Where("type_port = ?", info.TypePort)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return sysOperationRecords, total, err
}
