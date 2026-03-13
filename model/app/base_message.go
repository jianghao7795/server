/*
 * @Author: jianghao
 * @Date: 2022-10-17 10:05:23
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 11:39:11
 */
package app

import global "server/model"

type BaseMessage struct {
	global.MODEL
	Title        string `json:"title" form:"title" gorm:"column:title;comment:网站标题;size:50;"`
	Introduction string `json:"introduction" form:"introduction" gorm:"column:introduction;comment:网站简介;size:255"`
	HeadImg      string `json:"head_img" form:"head_img" query:"head_img" gorm:"column:head_img;comment:网站头像;size:255;"`
	Copyright    string `json:"copyright" form:"copyright" gorm:"column:copyright;comment:版权信息;size:255"`
	Link         string `json:"link" form:"link" gorm:"column:link;comment:链接;size:50;"`
	RecordInfo   string `json:"record_info" form:"record_info" query:"record_info" gorm:"column:record_info;comment:备案信息;size:255"`
	UserId       uint   `json:"user_id" form:"user_id" query:"user_id" gorm:"column:user_id;comment:用户"`
}

func (BaseMessage) TableName() string {
	return "base_messages"
}
