package frontend

import (
	"errors"
	global "server/model"
	"server/model/common/response"
	"server/model/frontend"
	loginRequest "server/model/frontend/request"
	"server/model/system"
	"server/utils"

	systemReq "server/model/system/request"
	systemRes "server/model/system/response"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type User struct{}

// var store = base64Captcha.DefaultMemStore

// var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
// var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
// func (u *FrontendUser) Login(c *fiber.Ctx) error {
// 	var user loginRequest.LoginForm
// 	_ = c.QueryParser(&user)
// 	if err := utils.Verify(user, utils.LoginVerifyFrontend); err != nil {
// 		return response.FailWithMessage(err.Error(), c)
// 		return
// 	}
// 	userInfo, err := userServiceApp.Login(user)
// 	if err != nil {
// 		global.LOG.Error(err.Error(), zap.Error(err))
// 		return response.FailWithMessage(err.Error(), c)
// 	} else {
// 		return response.OkWithDetailed(userInfo, "登录成功", c)
// 	}
// }

// Login 前台用户登录
// @Tags Frontend User
// @Summary 前台用户登录
// @Description 前台用户登录获取 token
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.Login true "登录信息"
// @Success 200 {object} response.Response{msg=string} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "登录失败"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /frontend/login [post]
func (b *User) Login(c *fiber.Ctx) error {
	var l systemReq.Login
	if err := c.BodyParser(&l); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := utils.Verify(l, utils.LoginVerifyFrontend); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if user, err := userService.Login(l.Username, l.Password); err != nil {
		global.LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		return response.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		return b.tokenNext(c, *user)
	}
	// if store.Verify(l.CaptchaId, l.Captcha, true) {

	// } else {
	// 	return response.FailWithMessage("验证码错误", c)
	// }
}

// 登录以后签发jwt
func (u *User) tokenNext(c *fiber.Ctx, user system.SysUser) error {
	j := utils.NewJWT() // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.LOG.Error("获取token失败!", zap.Error(err))
		return response.FailWithMessage("获取token失败", c)
	}
	c.Locals("frontend_user", user)
	if !global.CONFIG.System.UseMultipoint {
		return response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
		}, "登录成功", c)
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.LOG.Error("设置登录状态失败!", zap.Error(err))
			return response.FailWithMessage("设置登录状态失败", c)
		}
		return response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
		}, "登录成功", c)
	} else if err != nil {
		global.LOG.Error("设置登录状态失败!", zap.Error(err))
		return response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			return response.FailWithMessage("jwt作废失败", c)
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			return response.FailWithMessage("设置登录状态失败", c)
		}
		return response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
		}, "登录成功", c)
	}
}

func (*User) RegisterUser(c *fiber.Ctx) error {
	var userInfo loginRequest.RegisterUser
	_ = c.BodyParser(&userInfo)
	if userInfo.Password != userInfo.RePassword {
		global.LOG.Error("密码不一致!", zap.Error(errors.New("密码不一致")))
		return response.FailWithMessage("密码不一致", c)
	}
	if err := utils.Verify(userInfo, utils.RegisterVerifyFrontend); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}

	err := userServiceApp.RegisterUser(&userInfo)
	if err != nil {
		global.LOG.Error("注册失败!", zap.Error(err))
		return response.FailWithDetailed(nil, err.Error(), c)
	} else {
		return response.OkWithDetailed(nil, "注册成功", c)
	}
}

func (u *User) GetCurrent(c *fiber.Ctx) error {
	uuid := utils.GetUserUuid(c)
	if ReqUser, err := userService.GetUserInfo(uuid); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(ReqUser, "获取成功", c)
	}
}

func (u *User) UpdatePassword(c *fiber.Ctx) error {
	var resetPassword frontend.ResetPassword
	if err := c.BodyParser(&resetPassword); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}

	if err := utils.Verify(resetPassword, utils.ResetPasswordVerifyFrontend); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}

	if resetPassword.NewPassword != resetPassword.RepeatNewPassword {
		global.LOG.Error("密码不一致!", zap.Error(errors.New("密码不一致")))
		return response.FailWithMessage("密码不一致", c)
	}
	// resetPassword.Password = utils.MD5V([]byte(resetPassword.Password))
	// resetPassword.NewPassword = utils.MD5V([]byte(resetPassword.NewPassword))
	// resetPassword.RepeatNewPassword = utils.MD5V([]byte(resetPassword.RepeatNewPassword))
	if err := userServiceApp.ResetPassword(&resetPassword); err != nil {
		return response.FailWithMessage("重置密码失败："+err.Error(), c)
	}

	return response.OkWithDetailed(nil, "重置密码成功", c)
}

func (u *User) UpdateUserBackgroudImage(c *fiber.Ctx) error {
	var user frontend.User
	var err error
	err = c.BodyParser(&user)
	if err != nil {
		return response.FailWithMessage("获取数据失败", c)
	}
	err = userServiceApp.UpdateUserBackgroudImage(&user)
	if err != nil {
		return response.FailWithMessage("更新失败："+err.Error(), c)
	}
	return response.OkWithDetailed(nil, "更新成功", c)
}

func (u *User) UpdateUser(c *fiber.Ctx) error {
	var user frontend.User
	var err error
	err = c.BodyParser(&user)
	if err != nil {
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = utils.Verify(user, utils.UpdateUserVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if err = userServiceApp.UpdateUser(&user); err != nil {
		return response.FailWithDetailed(err.Error(), "更新失败", c)
	}
	return response.OkWithDetailed(nil, "更新成功", c)
}
