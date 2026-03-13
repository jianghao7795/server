package mobile

import (
	"errors"
	global "server/model"
	"server/model/mobile"
	"server/model/mobile/request"
	"server/model/mobile/response"
	"server/utils"

	"gorm.io/gorm"
)

type MobileLoginService struct{}

func (*MobileLoginService) Login(data *mobile.Login) (m response.LoginResponse, err error) {
	var user mobile.MobileUser
	// data.Password = utils.MD5V([]byte(data.Password))
	err = global.DB.Where("username = ? and password = ?", data.Username, data.Password).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return m, errors.New("密码错误")
	}

	j := utils.NewJWT()

	tokenString, expiresAt, err := j.MakeToken(*data, user.ID)
	if err != nil {
		return
	}
	m.Token = tokenString
	m.ExpiresAt = expiresAt
	return
}

func (*MobileLoginService) GetUserInfo(id uint) (u mobile.MobileUser, err error) {
	var user mobile.MobileUser
	err = global.DB.Model(&mobile.MobileUser{}).Where("id = ?", id).First(&user).Error
	return user, err
}

func (*MobileLoginService) UpdateUser(data *request.MobileUpdate, id uint) (err error) {
	return global.DB.Model(&mobile.MobileUser{}).Where("id = ?", id).Update(data.Field, data.Value).Error
}

func (*MobileLoginService) UpdatePassword(data request.MobileUpdatePassword) (err error) {
	var user mobile.MobileUser
	// passwordMoken := utils.MD5V([]byte(data.Password))
	db := global.DB.Model(&mobile.MobileUser{})
	err = db.Where("id = ? and password = ?", data.ID, data.Password).First(&user).Error
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("未找到用户，请更新请求")
	}
	err = db.Update("password", data.NewPassword).Error
	return err
}
