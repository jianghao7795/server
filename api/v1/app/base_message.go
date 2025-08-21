package app

import (
	"errors"
	"server-fiber/model/app"
	"server-fiber/model/common/response"

	global "server-fiber/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateBaseMessage 创建基础消息
// @Tags BaseMessage
// @Summary 创建基础消息
// @Description 创建新的基础消息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.BaseMessage true "基础消息信息"
// @Success 200 {object} response.Response{msg=string,data=integer,code=integer} "创建基础消息成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /base_message/createBaseMessage [post]
func (a *BaseMessageApi) CreateBaseMessage(c *fiber.Ctx) error {
	var baseMessage app.BaseMessage
	err := c.BodyParser(&baseMessage)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := baseMessageService.CreateBaseMessage(&baseMessage); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithId("创建成功", baseMessage.ID, c)
	}
}

// UpdateBaseMessage 更新基础消息
// @Tags BaseMessage
// @Summary 更新基础消息
// @Description 根据ID更新基础消息信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "基础消息ID" minimum(1)
// @Param data body app.BaseMessage true "基础消息信息"
// @Success 200 {object} response.Response{msg=string,code=integer} "更新基础消息成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /base_message/updateBaseMessage/{id} [put]
func (a *BaseMessageApi) UpdateBaseMessage(c *fiber.Ctx) error {
	var baseMessage app.BaseMessage
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	err = c.BodyParser(&baseMessage)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = baseMessageService.UpdateBaseMessage(id, &baseMessage); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// FindBaseMessage 查找基础消息
// @Tags BaseMessage
// @Summary 查找基础消息
// @Description 根据ID查找基础消息详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "基础消息ID" minimum(1)
// @Success 200 {object} response.Response{msg=string,data=app.BaseMessage,code=integer} "查找基础消息成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "基础消息不存在"
// @Router /base_message/getBaseMessage/{id} [get]
func (a *BaseMessageApi) FindBaseMessage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败!", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if responseBaseMessage, err := baseMessageService.FindBaseMessage(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// respBaseMessage := baseMessageNotFound{message: "not found"}
			str := "base message not found"
			global.LOG.Error("查询失败!", zap.Error(errors.New(str)))
			return response.OkWithData(str, c)
		} else {
			global.LOG.Error("查询失败!", zap.Error(err))
			return response.FailWithMessage("查询失败", c)
		}
	} else {
		return response.OkWithDetailed(responseBaseMessage, "查询成功", c)
	}
}
