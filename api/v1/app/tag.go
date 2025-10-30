package app

import (
	global "server-fiber/model"
	"server-fiber/model/app"
	appReq "server-fiber/model/app/request"
	"server-fiber/model/common/request"
	"server-fiber/model/common/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CreateTag 创建标签
// @Tags Tag
// @Summary 创建标签
// @Description 创建新的标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.Tag true "标签信息"
// @Success 200 {object} response.Response{msg=string} "创建标签成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/createTag [post]
func (TagApi *TagApi) CreateTag(c *fiber.Ctx) error {
	var appTab app.Tag
	err := c.BodyParser(&appTab)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := appTabService.CreateTag(&appTab); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败", c)
	} else {
		return response.OkWithId("创建成功", appTab.ID, c)
	}
}

// DeleteTag 删除标签
// @Tags Tag
// @Summary 删除标签
// @Description 根据标签ID删除指定标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "标签ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除标签成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/deleteTag/{id} [delete]
func (TagApi *TagApi) DeleteTag(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.FailWithDetailed(fiber.Map{"msg": "获取id参数失败"}, "参数错误", c)
	}
	if err := appTabService.DeleteTag(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithMessage("删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// DeleteTagByIds 批量删除标签
// @Tags Tag
// @Summary 批量删除标签
// @Description 根据ID列表批量删除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "标签ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除标签成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/deleteTagByIds [delete]
func (TagApi *TagApi) DeleteTagByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	_ = c.BodyParser(&IDS)
	if err := appTabService.DeleteTagByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithMessage("批量删除失败", c)
	} else {
		return response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateTag 更新标签
// @Tags Tag
// @Summary 更新标签
// @Description 更新标签信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.Tag true "标签信息"
// @Success 200 {object} response.Response{msg=string} "更新标签成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/updateTag [put]
func (TagApi *TagApi) UpdateTag(c *fiber.Ctx) error {
	var appTab app.Tag
	err := c.BodyParser(&appTab)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = appTabService.UpdateTag(&appTab); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败", c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// FindTag 用id查询Tag
// @Tags Tag
// @Summary 用id查询Tag
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query app.Tag true "用id查询Tag"
// @Success 200 {string} string "{"success":true,"data":app.Tag,"msg":"查询成功"}"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/:id [get]
func (TagApi *TagApi) FindTag(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if tag, err := appTabService.GetTag(uint(id)); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败", c)
	} else {
		return response.OkWithData(tag, c)
	}
}

// GetTagList 分页获取Tag列表
// @Tags Tag
// @Summary 分页获取Tag列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appReq.TagSearch true "分页获取Tag列表"
// @Success 200 {string} string "{"success":true,"data":[]app.Tag,total=integer,"msg":"获取成功"}"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /tag/getTagList [get]
func (TagApi *TagApi) GetTagList(c *fiber.Ctx) error {
	var pageInfo appReq.TagSearch
	_ = c.QueryParser(&pageInfo)
	if list, total, err := appTabService.GetTagInfoList(&pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
