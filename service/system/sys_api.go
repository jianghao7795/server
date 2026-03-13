package system

import (
	"errors"
	"fmt"
	global "server/model"
	"server/model/common/request"
	"server/model/system"
	systemReq "server/model/system/request"

	"gorm.io/gorm"
)

//
//@function: CreateApi
//@description: 新增基础api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) CreateApi(api *system.SysApi) (err error) {
	if !errors.Is(global.DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return global.DB.Create(api).Error
}

//
//@function: DeleteApi
//@description: 删除基础api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) DeleteApi(api system.SysApi) (err error) {
	err = global.DB.Delete(&api).Error
	CasbinServiceApp.ClearCasbin(1, api.Path, api.Method)
	return err
}

//
//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: api model.SysApi, info request.PageInfo, order string, desc bool
//@return: err error

func (apiService *ApiService) GetAPIInfoList(info *systemReq.SearchApiParams) (list []system.SysApi, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&system.SysApi{})
	if info.Path != "" {
		db = db.Where("path like ?", "%"+info.Path+"%")
	}
	if info.Description != "" {
		db = db.Where("description like ?", "%"+info.Description+"%")
	}
	if info.ApiGroup != "" {
		db = db.Where("api_group like ?", "%"+info.ApiGroup+"%")
	}
	if info.Method != "" {
		db = db.Where("api_group = ?", info.Method)
	}
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	sort := "id"
	if info.OrderKey != "" {
		if info.Desc == "true" {
			sort = fmt.Sprintf("%s %s", info.OrderKey, "desc")
		} else {
			sort = info.OrderKey
		}
	}
	err = db.Limit(limit).Offset(offset).Order(sort).Find(&list).Error
	return
}

//
//@function: GetAllApis
//@description: 获取所有的api
//@return: err error, apis []model.SysApi

func (apiService *ApiService) GetAllApis() (apis []system.SysApi, err error) {
	err = global.DB.Find(&apis).Error
	return
}

//
//@function: GetApiById
//@description: 根据id获取api
//@param: id float64
//@return: err error, api model.SysApi

func (apiService *ApiService) GetApiById(id int) (api system.SysApi, err error) {
	err = global.DB.Where("id = ?", id).First(&api).Error
	return
}

//
//@function: UpdateApi
//@description: 根据id更新api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) UpdateApi(api *system.SysApi) (err error) {
	var oldA system.SysApi
	err = global.DB.Where("id = ?", api.ID).First(&oldA).Error
	if err != nil {
		return err
	}
	if oldA.Path != api.Path || oldA.Method != api.Method {
		if !errors.Is(global.DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同api路径")
		}
	}
	return global.DB.Save(api).Error
}

//
//@function: DeleteApis
//@description: 删除选中API
//@param: apis []model.SysApi
//@return: err error

func (apiService *ApiService) DeleteApisByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]system.SysApi{}, "id in ?", ids.Ids).Error
	return err
}

func (apiService *ApiService) DeleteApiByIds(ids []string) (err error) {
	return global.DB.Delete(&system.SysApi{}, "id in ?", ids).Error
}
