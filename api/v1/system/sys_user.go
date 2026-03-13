package system

import (
	"errors"
	global "server/model"
	"server/model/common/response"
	"server/model/system"
	"server/utils"
	"strconv"

	systemReq "server/model/system/request"
	systemRes "server/model/system/response"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Login 用户登录
// @Tags Base
// @Summary 用户登录
// @Description 用户登录获取 token 和用户信息
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.Login true "登录信息"
// @Success 200 {object} response.Response{msg=string} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "登录失败"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /base/login [post]
func (b *BaseApi) Login(c *fiber.Ctx) error {
	var l systemReq.Login
	if err := c.BodyParser(&l); err != nil {
		global.LOG.Error("参数解析失败!", zap.Error(err))
		return response.FailWithMessage("参数错误", c)
	}
	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if store.Verify(l.CaptchaId, l.Captcha, true) {
		if user, err := userService.Login(l.Username, l.Password); err != nil {
			global.LOG.Error(err.Error(), zap.Error(err))
			errorMessage := err.Error()
			if err == gorm.ErrRecordNotFound {
				errorMessage = "账户或密码错误"
			}
			return response.FailWithMessage("登录失败："+errorMessage, c)
		} else {
			return b.tokenNext(c, user)
		}
	} else {
		return response.FailWithMessage("验证码错误", c)
	}
}

// LoginToken 用户登录获取token
// @Tags Base
// @Summary 用户登录获取token
// @Description 用户登录获取 token，无需验证码
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.LoginToken true "登录信息"
// @Success 200 {object} response.Response{msg=string} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "登录失败"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /base/getToken/login [post]
func (b *BaseApi) LoginToken(c *fiber.Ctx) error {
	var l systemReq.LoginToken
	_ = c.BodyParser(&l)
	if err := utils.Verify(l, utils.TokenLoginVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	u := &system.SysUser{Username: l.Username, Password: l.Password}
	if user, err := userService.LoginToken(u); err != nil {
		global.LOG.Error(err.Error(), zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		return b.tokenNext(c, user)
	}
}

// 登录以后签发jwt
func (b *BaseApi) tokenNext(c *fiber.Ctx, user *system.SysUser) error {
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
		global.LOG.Error(err.Error(), zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	if !global.CONFIG.System.UseMultipoint {
		return response.OkWithDetailed(systemRes.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
		}, "登录成功", c)
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); errors.Is(err, redis.Nil) {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.LOG.Error("设置登录状态失败!", zap.Error(err))
			return response.FailWithMessage("设置登录状态失败", c)
		}
		return response.OkWithDetailed(systemRes.LoginResponse{
			User:      *user,
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
			User:      *user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
		}, "登录成功", c)
	}
}

// Register
// @Tags SysUser
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body systemReq.Register true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {object} response.Response{msg=string} "用户注册账号,返回包括用户信息"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/admin_register [post]
func (b *BaseApi) Register(c *fiber.Ctx) error {
	var r systemReq.Register
	_ = c.BodyParser(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(user)
	if err != nil {
		global.LOG.Error("注册失败!", zap.Error(err))
		return response.FailWithDetailed(systemRes.SysUserResponse{User: *userReturn}, "注册失败", c)
	} else {
		return response.OkWithDetailed(systemRes.SysUserResponse{User: *userReturn}, "注册成功", c)
	}
}

// ChangePassword
// @Tags SysUser
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body systemReq.ChangePasswordStruct true "用户名, 原密码, 新密码"
// @Success 200 {object} response.Response{msg=string} "用户修改密码"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/changePassword [post]
func (b *BaseApi) ChangePassword(c *fiber.Ctx) error {
	var user systemReq.ChangePasswordStruct
	_ = c.BodyParser(&user)
	if err := utils.Verify(user, utils.ChangePasswordVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	u := &system.SysUser{Username: user.Username, Password: user.Password}
	if respUser, err := userService.ChangePassword(u, user.NewPassword); err != nil {
		global.LOG.Error(respUser.Username+"修改失败!", zap.Error(err))
		return response.FailWithMessage(respUser.Username+"修改失败，原密码与当前账户不符", c)
	} else {
		return response.OkWithMessage(respUser.Username+"修改成功", c)
	}
}

// GetUserList
// @Tags SysUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param query query request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=object,msg=string} "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/getUserList [get]
func (b *BaseApi) GetUserList(c *fiber.Ctx) error {
	var searchInfo systemReq.SearchInfo
	_ = c.QueryParser(&searchInfo)
	if err := utils.Verify(searchInfo, utils.PageInfoVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}
	if list, total, err := userService.GetUserInfoList(searchInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     searchInfo.Page,
			PageSize: searchInfo.PageSize,
		}, "获取成功", c)
	}
}

// SetUserAuthority 更改用户权限
// @Tags SysUser
// @Summary 更改用户权限
// @Description 更改指定用户的角色权限
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.SetUserAuth true "用户权限信息"
// @Success 200 {object} response.Response{msg=string} "设置用户权限成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *fiber.Ctx) error {
	var sua systemReq.SetUserAuth
	_ = c.BodyParser(&sua)
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		return response.FailWithMessage(UserVerifyErr.Error(), c)
	}
	userID, _ := utils.GetUserID(c)
	uuid := utils.GetUserUuid(c)
	if err := userService.SetUserAuthority(userID, uuid, sua.AuthorityId); err != nil {
		global.LOG.Error("修改失败!", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	} else {
		claims := utils.GetUserInfo(c)
		j := utils.NewJWT() // 唯一签名
		claims.AuthorityId = sua.AuthorityId
		if token, err := j.CreateToken(*claims); err != nil {
			global.LOG.Error("修改失败!", zap.Error(err))
			return response.FailWithMessage(err.Error(), c)
		} else {
			c.Set("new-token", token)
			c.Set("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
			return response.OkWithMessage("修改成功", c)
		}

	}
}

// SetUserAuthorities 设置用户权限
// @Tags SysUser
// @Summary 设置用户权限
// @Description 设置指定用户的多个角色权限
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.SetUserAuthorities true "用户权限信息"
// @Success 200 {object} response.Response{msg=string} "设置用户权限成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *fiber.Ctx) error {
	var sua systemReq.SetUserAuthorities
	_ = c.BodyParser(&sua)
	if err := userService.SetUserAuthorities(sua.ID, sua.AuthorityIds); err != nil {
		global.LOG.Error("修改失败!", zap.Error(err))
		return response.FailWithMessage("修改失败", c)
	} else {
		return response.OkWithMessage("修改成功", c)
	}
}

// DeleteUser 删除用户
// @Tags SysUser
// @Summary 删除用户
// @Description 根据用户ID删除指定用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "用户ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除用户成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "不能删除自己"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/deleteUser/{id} [delete]
func (b *BaseApi) DeleteUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	jwtId, _ := utils.GetUserID(c)
	if jwtId == uint(id) {
		return response.FailWithMessage("删除失败, 自杀失败", c)
	}
	if err := userService.DeleteUser(id); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// SetUserInfo 设置用户信息
// @Tags SysUser
// @Summary 设置用户信息
// @Description 设置指定用户的基本信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysUser true "用户信息"
// @Success 200 {object} response.Response{msg=string} "设置用户信息成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *fiber.Ctx) error {
	var user systemReq.ChangeUserInfo
	_ = c.BodyParser(&user)
	if err := utils.Verify(user, utils.IdVerify); err != nil {
		return response.FailWithMessage(err.Error(), c)
	}

	if len(user.AuthorityIds) != 0 {
		err := userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			global.LOG.Error("设置失败!", zap.Error(err))
			return response.FailWithMessage("设置失败", c)
		}
	}

	if err := userService.SetUserInfo(system.SysUser{
		MODEL: global.MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
	}); err != nil {
		global.LOG.Error("设置失败!", zap.Error(err))
		return response.FailWithMessage("设置失败", c)
	} else {
		return response.OkWithMessage("设置成功", c)
	}
}

// SetSelfInfo
// @Tags SysUser
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysUser true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {object} response.Response{msg=string} "设置用户信息"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/SetSelfInfo [put]
func (b *BaseApi) SetSelfInfo(c *fiber.Ctx) error {
	var user systemReq.ChangeUserInfo
	_ = c.BodyParser(&user)
	user.ID, _ = utils.GetUserID(c)
	if err := userService.SetUserInfo(system.SysUser{
		MODEL: global.MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
	}); err != nil {
		global.LOG.Error("设置失败!", zap.Error(err))
		return response.FailWithMessage("设置失败", c)
	} else {
		return response.OkWithMessage("设置成功", c)
	}
}

// GetUserInfo
// @Tags SysUser
// @Summary 获取用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取用户信息"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *fiber.Ctx) error {
	uuid := utils.GetUserUuid(c)
	if ReqUser, err := userService.GetUserInfo(uuid); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(ReqUser, "获取成功", c)
	}
}

// ResetPassword
// @Tags SysUser
// @Summary 重置用户密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body system.SysUser true "ID"
// @Success 200 {object} response.Response{msg=string} "重置用户密码"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/resetPassword [post]
func (b *BaseApi) ResetPassword(c *fiber.Ctx) error {
	var user system.SysUser
	_ = c.BodyParser(&user)
	if err := userService.ResetPassword(user.ID); err != nil {
		global.LOG.Error("重置失败!", zap.Error(err))
		return response.FailWithMessage("重置失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("重置成功", c)
	}
}

// GetUserCount
// @Tags SysUser
// @Summary 获取人员总数
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取人员总数"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /user/getUserCount [get]
func (b *BaseApi) GetUserCount(c *fiber.Ctx) error {
	if userCount, err := userService.UserCount(); err != nil {
		global.LOG.Error("获取总数失败!", zap.Error(err))
		return response.FailWithMessage("获取总数失败"+err.Error(), c)
	} else {
		return response.OkWithDetailed(fiber.Map{"count": userCount}, "获取成功", c)
	}
}

func (b *BaseApi) GetFlow(c *fiber.Ctx) error {
	receiveBytes, transmitBytes, _ := utils.TotalFlowByDevice("lo")
	return response.OkWithDetailed(fiber.Map{"receiveBytes": receiveBytes, "transmitBytes": transmitBytes}, "获取成功", c)
}
