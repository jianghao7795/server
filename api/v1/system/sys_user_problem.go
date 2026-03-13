package system

import (
	"strconv"

	global "server/model"
	"server/model/common/response"
	problemReq "server/model/system"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// UserProblem 用户问题管理
type UserProblem struct{}

// GetProblemSetting 获取用户问题设置
// @Summary 获取用户问题设置
// @Description 根据用户ID获取用户的问题设置列表
// @Tags 用户问题管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]problemReq.SysUserProblem},msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "传错参数,请传user id"
// @Failure 500 {object} response.Response{msg=string} "获取失败"
// @Router /api/v1/system/user-problem/setting/{id} [get]
func (*UserProblem) GetProblemSetting(c fiber.Ctx) error {
	var search problemReq.SysUserProblem
	var err error
	search.SysUserId, err = strconv.Atoi(c.Params("id"))
	if err != nil {
		global.LOG.Error("传参错误!", zap.Error(err))
		return response.FailWithMessage("传错参数,请传user id", c)
	}
	list, err := userProblem.GetUserProblemSettingList(&search)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	}
	return response.OkWithDetailed(response.PageResult{
		List: list,
	}, "获取成功", c)
}

// UpdateProblemSettingData 更新问题设置数据结构
type UpdateProblemSettingData struct {
	Data []problemReq.SysUserProblem `json:"data" example:"[{\"question\":\"安全问题1\",\"answer\":\"答案1\"}]"`
}

// UpdateProblemSetting 更新用户问题设置
// @Summary 更新用户问题设置
// @Description 批量更新用户的安全问题设置
// @Tags 用户问题管理
// @Accept json
// @Produce json
// @Param data body UpdateProblemSettingData true "问题设置数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Failure 400 {object} response.Response{msg=string} "传错参数"
// @Router /api/v1/system/user-problem/setting [put]
func (*UserProblem) UpdateProblemSetting(c fiber.Ctx) error {
	var dataProblem UpdateProblemSettingData
	err := c.Bind().Body(&dataProblem)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			global.LOG.Error("传参错误!", zap.Error(err))
			return response.FailWithDetailed(err.Error(), "传错参数", c)
		}
		global.LOG.Error("传参错误!", zap.Error(err))

		return response.FailWithDetailed(errs.Translate(global.Validate), "传错参数", c)
	}
	message, err := userProblem.SetUserProblemSetting(dataProblem.Data)
	if err != nil {
		global.LOG.Error("传参错误!", zap.Error(err))
		return response.FailWithDetailed(err, "传错参数", c)
	}
	return response.OkWithMessage(message, c)
}

// HasSetting 检查用户是否已设置问题
// @Summary 检查用户是否已设置问题
// @Description 检查指定用户是否已经设置了安全问题
// @Tags 用户问题管理
// @Accept json
// @Produce json
// @Param uid path int true "用户ID"
// @Success 200 {object} response.Response{data=bool,msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "传错参数"
// @Router /api/v1/system/user-problem/has-setting/{uid} [get]
func (*UserProblem) HasSetting(c fiber.Ctx) error {
	SysUserProblemId, _ := strconv.Atoi(c.Params("uid"))
	isSetting, err := userProblem.HasSetting(SysUserProblemId)
	if err != nil {
		global.LOG.Error("传参错误!", zap.Error(err))
		return response.FailWithDetailed(err, "传错参数", c)
	}
	return response.OkWithDetailed(isSetting, "获取成功", c)
}

// VerifyProblemSettingData 验证问题答案数据结构
type VerifyProblemSettingData struct {
	Data problemReq.SysUserProblem `json:"data" example:"{\"question\":\"安全问题1\",\"answer\":\"答案1\"}"`
}

// VerifyAnswer 验证问题答案
// @Summary 验证问题答案
// @Description 验证用户提供的安全问题答案是否正确
// @Tags 用户问题管理
// @Accept json
// @Produce json
// @Param data body VerifyProblemSettingData true "问题验证数据"
// @Success 200 {object} response.Response{data=bool,msg=string} "已验证"
// @Failure 400 {object} response.Response{msg=string} "传错参数"
// @Failure 404 {object} response.Response{msg=string} "未查到此问题"
// @Router /api/v1/system/user-problem/verify [post]
func (*UserProblem) VerifyAnswer(c fiber.Ctx) error {
	var dataProblem VerifyProblemSettingData
	err := c.Bind().Body(&dataProblem)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			global.LOG.Error("传参错误!", zap.Error(err))
			return response.FailWithDetailed(err.Error(), "传错参数", c)
		}
		global.LOG.Error("传参错误!", zap.Error(err))
		return response.FailWithDetailed(errs.Translate(global.Validate), "传错参数", c)
	}
	ispassed, err := userProblem.VerifyAnswer(&dataProblem.Data)
	if err != nil {
		global.LOG.Error("未查到此问题!", zap.Error(err))
		return response.FailWithDetailed(err, "未查到此问题", c)
	}
	return response.OkWithDetailed(ispassed, "已验证", c)
}
