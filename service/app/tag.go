package app

import (
	global "server/model"
	"server/model/app"
	appReq "server/model/app/request"
	"server/model/common/request"
	"strings"
)

type TagService struct{}

// CreateTag 创建Tag记录
// Author wuhao
func (tagService *TagService) CreateTag(tag *app.Tag) (err error) {
	err = global.DB.Create(tag).Error
	return err
}

// DeleteTag 删除Tag记录
// Author wuhao
func (tagService *TagService) DeleteTag(id uint) (err error) {
	err = global.DB.Delete(&app.Tag{}, id).Error
	return err
}

// DeleteTagByIds 批量删除Tag记录
// Author wuhao
func (tagService *TagService) DeleteTagByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]app.Tag{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateTag 更新Tag记录
// Author wuhao
func (tagService *TagService) UpdateTag(tag *app.Tag) (err error) {
	var tagNew app.Tag
	// var count int64
	err = global.DB.Where("id = ?", tag.ID).First(&tagNew).Error
	if err != nil {
		return err
	}
	err = global.DB.Save(tag).Error
	return err
}

// GetTag 根据id获取Tag记录
// Author wuhao
func (tagService *TagService) GetTag(id uint) (tag app.Tag, err error) {
	err = global.DB.Where("id = ?", id).First(&tag).Error
	return
}

// GetTagInfoList 分页获取Tag记录
// Author wuhao
func (tagService *TagService) GetTagInfoList(info *appReq.TagSearch) (list []app.Tag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&app.Tag{})
	var tags []app.Tag
	if info.Name != "" {
		db = db.Where("name like ?", strings.Join([]string{"%", info.Name, "%"}, ""))
	}
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&tags).Error
	return tags, total, err
}
