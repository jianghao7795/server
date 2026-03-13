package mobile

import (
	global "server/model"
	"server/model/common/response"
	"server/model/mobile"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type RegisterMobile struct{}

// Register 移动端用户注册
// @Tags Mobile Register
// @Summary 移动端用户注册
// @Description 移动端用户注册新账号
// @Accept application/json
// @Produce application/json
// @Param data body mobile.Register true "注册信息"
// @Success 200 {object} response.Response{msg=string} "注册成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /mobile/register [post]
func (*RegisterMobile) Register(c fiber.Ctx) (err error) {
	var data mobile.Register
	err = c.Bind().Body(&data)
	if err != nil {
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = registerService.Register(data); err != nil {
		global.LOG.Error("注册失败!", zap.Error(err))
		return response.FailWithMessage400("注册失败，请重试", c)
	} else {
		return response.OkWithDetailed("", "注册成功", c)
	}
}
