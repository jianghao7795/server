package app

import (
	"errors"
	global "server/model"
	"server/model/app"
)

type BaseMessageService struct{}

/**
 * @description: 创建baseMessage
 * @param {app.BaseMessage} baseMessage
 * @return {error}
 */
func (*BaseMessageService) CreateBaseMessage(baseMessage *app.BaseMessage) (err error) {
	err = global.DB.Create(baseMessage).Error
	return
}

/**
 * @description: 更新baseMessage
 * @param {int, app.BaseMessage} id, baseMessage
 * @return {error}
 */
func (*BaseMessageService) UpdateBaseMessage(id int, baseMessage *app.BaseMessage) (err error) {
	var baseMessageReplica app.BaseMessage
	db := global.DB.Model(&app.BaseMessage{}).Where("id = ?", id).First(&baseMessageReplica)
	if baseMessageReplica.ID == 0 {
		return errors.New("数据库没有记录")
	}
	// 这里可以添加对 baseMessageReplica 的字段赋值
	result := db.Save(baseMessage)
	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected == 0 {
		err = errors.New("没有更新任何数据")
		return
	}
	return
}

/**
 * @description: 获取baseMessage
 * @param {uint} id
 * @return {app.BaseMessage, error}
 */
func (*BaseMessageService) FindBaseMessage(id uint) (app.BaseMessage, error) {
	var base app.BaseMessage
	err := global.DB.Where("user_id = ?", id).First(&base).Error
	return base, err
}
