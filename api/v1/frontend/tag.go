package frontend

import (
	"strconv"

	global "server/model"
	appReq "server/model/app/request"
	"server/model/common/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type TagApi struct{}

// GetTagList 获取标签列表
// @Tags Frontend Tag
// @Summary 获取标签列表
// @Description 分页获取前台展示的标签列表
// @Produce application/json
// @Param page query integer false "页码" default(1) minimum(1)
// @Param pageSize query integer false "每页数量" default(10) minimum(1) maximum(100)
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]app.Tag},code=integer} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Router /frontend/getTagList [get]
func (appTabApi *TagApi) GetTagList(c fiber.Ctx) error {
	var pageInfo appReq.TagSearch
	_ = c.Bind().Query(&pageInfo)
	if list, err := tagServiceApp.GetTagList(c); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}

// GetTag 获取标签详情
// @Tags Frontend Tag
// @Summary 获取标签详情
// @Description 根据标签ID获取标签及其关联的文章
// @Produce application/json
// @Param id path integer true "标签ID" minimum(1)
// @Success 200 {object} response.Response{msg=string,data=object{tag=[]app.Tag},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /frontend/getTag/{id} [get]
func (appTabApi *TagApi) GetTag(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取Ids失败", c)
	}
	if tagArticles, err := tagServiceApp.GetTagArticle(id, c); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithData(fiber.Map{"tag": tagArticles}, c)
	}
}
