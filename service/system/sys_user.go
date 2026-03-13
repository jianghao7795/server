package system

import (
	"errors"
	"fmt"
	global "server/model"
	"server/model/system"
	"server/model/system/request"
	"server/utils"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: err error, userInter model.SysUser

func (userService *UserService) Register(u *system.SysUser) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码sha512简单加密 注册
	u.Password = utils.Sha512V(u.Password)
	u.UUID = uuid.New()
	// u.ExpirationDate = time.Now().Format("2006-01-02 15:04:05")
	err = global.DB.Create(&u).Error
	return u, err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter g*model.SysUser

func (userService *UserService) Login(username, password string) (*system.SysUser, error) {
	if nil == global.DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	// u.Password = utils.MD5V([]byte(u.Password))
	err := global.DB.Where("username = ? AND password = ?", username, utils.Sha512V(password)).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		var am system.SysMenu
		ferr := global.DB.First(&am, "name = ? AND authority_id = ?", user.Authority.DefaultRouter, user.AuthorityId).Error
		if errors.Is(ferr, gorm.ErrRecordNotFound) {
			user.Authority.DefaultRouter = "404"
		}
		if user.HeaderImg == "" {
			user.HeaderImg = "public/logo.png"
		} else {
			if isNotExit, _ := utils.PathExistsFile(user.HeaderImg); !isNotExit {
				user.HeaderImg = "public/logo.png"
			}
		}

	} else {
		return nil, err
	}

	// if user.UpdatedAt.AddDate(0, 3, 0).After(time.Now()) {
	// 	return nil, errors.New("密码已过期，请修改密码")
	// }

	return &user, nil
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: LoginToken md5 加密
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) LoginToken(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	// log.Println("user: ", u.Username, u.Password)
	// u.Password = utils.MD5V([]byte(u.Password))
	err = global.DB.Where("username = ? AND password = ?", u.Username, utils.Sha512V(u.Password)).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		var am system.SysMenu
		ferry := global.DB.First(&am, "name = ? AND authority_id = ?", user.Authority.DefaultRouter, user.AuthorityId).Error
		if errors.Is(ferry, gorm.ErrRecordNotFound) {
			user.Authority.DefaultRouter = "404"
		}
	}
	return &user, err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: err error, userInter *model.SysUser

func (userService *UserService) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	// u.Password = utils.MD5V([]byte(u.Password))
	err = global.DB.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).Update("password", newPassword).Error
	return u, err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list any, total int64

func (userService *UserService) GetUserInfoList(info request.SearchInfo) (list any, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Username != "" {
		db = db.Where("username = ?", info.Username)
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserService) SetUserAuthority(id uint, uuid uuid.UUID, authorityId string) (err error) {
	assignErr := global.DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUseAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = global.DB.Where("uuid = ?", uuid).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (userService *UserService) SetUserAuthorities(id uint, authorityIds []string) (err error) {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Delete(&[]system.SysUseAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		useAuthority := []system.SysUseAuthority{}
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUseAuthority{
				SysUserId:               id,
				SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (userService *UserService) DeleteUser(id int) (err error) {
	var user system.SysUser
	err = global.DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	err = global.DB.Delete(&[]system.SysUseAuthority{}, "sys_user_id = ?", id).Error
	return err
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	return global.DB.Updates(&req).Error
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	var am system.SysMenu
	ferr := global.DB.First(&am, "name = ? AND authority_id = ?", reqUser.Authority.DefaultRouter, reqUser.AuthorityId).Error
	if errors.Is(ferr, gorm.ErrRecordNotFound) {
		reqUser.Authority.DefaultRouter = "404"
	}
	return reqUser, err
}

//@author: wuhao
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.DB.Where("`id` = ?", id).First(&u).Error
	return &u, err
}

//@author: wuhao
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@author: [jianghao](https://github.com/JiangHaoCode)
//@function: resetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (userService *UserService) ResetPassword(ID uint) (err error) {
	err = global.DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.Sha512V("123456")).Error
	return err
}

func (*UserService) UserCount() (userCount int64, err error) {
	err = global.DB.Model(&system.SysUser{}).Count(&userCount).Error
	return
}
