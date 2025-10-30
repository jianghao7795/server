package frontend

import (
	"errors"
	global "server-fiber/model"
	"server-fiber/model/common/response"
	"server-fiber/model/frontend/request"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleApi struct{}

// GetArticleList 分页获取文章列表
// @Tags Frontend Article
// @Summary 分页获取文章列表
// @Description 分页获取前台文章列表，支持搜索和筛选
// @Produce application/json
// @Param page query integer false "页码" default(1) minimum(1)
// @Param pageSize query integer false "每页数量" default(10) minimum(1) maximum(100)
// @Param title query string false "文章标题搜索"
// @Param state query integer false "文章状态"
// @Param is_important query integer false "是否首页显示"
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]app.Article,total=integer,page=integer,pageSize=integer},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Router /getArticleList [get]
func (s *ArticleApi) GetArticleList(c *fiber.Ctx) error {
	var pageInfo request.ArticleSearch
	err := c.QueryParser(&pageInfo)
	if err != nil {
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "参数错误", c)
	}
	// log.Println("is_promint: ", pageInfo.IsImportant)
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}

	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}

	if list, total, err := articleServiceApp.GetArticleList(&pageInfo, c); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetArticleDetail 获取文章详情
// @Tags Frontend Article
// @Summary 获取文章详情
// @Description 根据文章ID获取文章详细信息
// @Produce application/json
// @Param id path integer true "文章ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /getArticle/{id} [get]
func (s *ArticleApi) GetArticleDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取Id失败!", zap.Error(err))
		return response.FailWithMessage("获取Id失败", c)

	}
	articleDetail, err := articleServiceApp.GetArticleDetail(id, c)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.FailWithMessage("文章没有，请重新查询", c)
	}
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取失败", c)

	} else {
		return response.OkWithData(articleDetail, c)
	}
}

// GetSearchArticle 搜索文章
// @Tags Frontend Article
// @Summary 搜索文章
// @Description 根据标签或文章名称搜索相关内容
// @Produce application/json
// @Param name path string true "搜索类型" Enums(tags, articles)
// @Param value path string true "搜索值"
// @Param sort query string false "排序方式"
// @Success 200 {object} response.Response{msg=string} "搜索成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /getSearchArticle/{name}/{value} [get]
func (s *ArticleApi) GetSearchArticle(c *fiber.Ctx) error {
	var searchValue request.ArticleSearch
	err := c.ParamsParser(&searchValue)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取数据失败", c)
	}
	searchValue.Sort = c.Query("sort")
	if searchValue.Name != "tags" && searchValue.Name != "articles" {
		return response.FailWithDetailed(fiber.Map{"msg": "查询的不是tag 或 article"}, "查询的不是tag 或 article", c)
	}
	if list, err := articleServiceApp.GetSearchArticle(&searchValue); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List: list,
			// Total:    total,
			// Page:     pageInfo.Page,
			// PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
