package app

import (
	"errors"
	"server/model/app"
	"server/model/common/request"

	global "server/model"

	appReq "server/model/app/request"

	"gorm.io/gorm"
)

type UserService struct{}

// CreateUser 创建User记录
// Author wuhao
func (userService *UserService) CreateUser(user *app.User) (err error) {
	var userCurrent app.User
	if !errors.Is(global.DB.Where("name = ?", user.Name).First(&userCurrent).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册")
	}
	// user.Password = utils.MD5V([]byte(user.Password))
	err = global.DB.Create(user).Error
	return err
}

// DeleteUser 删除User记录
// Author wuhao
func (userService *UserService) DeleteUser(ID int) (err error) {
	err = global.DB.Where("id = ?", ID).Delete(&app.User{}).Error
	return err
}

// DeleteUserByIds 批量删除User记录
// Author wuhao
func (userService *UserService) DeleteUserByIds(ids request.IdsReq) (err error) {
	err = global.DB.Delete(&[]app.User{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateUser 更新User记录
// Author wuhao
func (userService *UserService) UpdateUser(user *app.User) (err error) {
	var userReplice app.User
	db := global.DB.Where("id = ?", user.ID).First(&userReplice)
	if userReplice.ID == 0 {
		return errors.New("更新错误， 未能找到该用户")
	}
	return db.Save(user).Error
}

// GetUser 根据id获取User记录
// Author wuhao
func (userService *UserService) GetUser(id uint) (user app.User, err error) {
	err = global.DB.Where("id = ?", id).First(&user).Error
	return
}

// GetUserInfoList 分页获取User记录
// Author wuhao
func (userService *UserService) GetUserInfoList(info *appReq.UserSearch) (list []app.User, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&app.User{})

	var users []app.User
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&users).Error
	return users, total, err
}

func (userService *UserService) FindIsUser(id int) (bool, error) {
	var user app.User
	err := global.DB.Where("id = ?", id).First(&user).Error
	return errors.Is(err, gorm.ErrRecordNotFound), err
}
