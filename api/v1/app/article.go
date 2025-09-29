package app

import (
	global "server-fiber/model"
	"server-fiber/model/app"
	appReq "server-fiber/model/app/request"
	"server-fiber/model/common/request"
	"server-fiber/model/common/response"
	"server-fiber/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CreateArticle 创建文章
// @Tags Article
// @Summary 创建文章
// @Description 创建新的文章
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.Article true "文章信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /article/createArticle [post]
func (a *ArticleApi) CreateArticle(c *fiber.Ctx) error {
	var article app.Article
	if err := c.BodyParser(&article); err != nil {
		return utils.ErrorHandlerInstance.HandleValidationError(c, "article data", err)
	}

	if err := articleService.CreateArticle(&article); err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "创建文章", err)
	}

	return response.OkWithId("创建成功", article.ID, c)
}

// DeleteArticle 删除文章
// @Tags Article
// @Summary 删除文章
// @Description 根据文章ID删除指定文章
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "文章ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /article/deleteArticle/{id} [delete]
func (*ArticleApi) DeleteArticle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if err := articleService.DeleteArticle(uint(id)); err != nil {
		global.LOG.Error("删除失败", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "删除失败", c)

	}
	return response.OkWithMessage("删除成功", c)
}

// DeleteArticleByIds 批量删除文章
// @Tags Article
// @Summary 批量删除文章
// @Description 根据ID列表批量删除文章
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "文章ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /article/deleteArticleByIds [delete]
func (a *ArticleApi) DeleteArticleByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	err := c.BodyParser(&IDS)
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if err := articleService.DeleteArticleByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "批量删除失败", c)

	}
	return response.OkWithMessage("批量删除成功", c)
}

// UpdateArticle 更新文章
// @Tags Article
// @Summary 更新文章
// @Description 根据文章ID更新文章信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "文章ID" minimum(1)
// @Param data body app.Article true "文章信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /article/updateArticle/{id} [put]
func (*ArticleApi) UpdateArticle(c *fiber.Ctx) error {
	var article app.Article
	err := c.BodyParser(&article)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败: "+err.Error(), c)
	}
	var id int
	id, err = c.ParamsInt("id")
	if err != nil {
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, err.Error(), c)
	}
	article.ID = uint(id)
	if err = articleService.UpdateArticle(&article); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "更新失败", c)

	}
	return response.OkWithMessage("更新成功", c)
}

// FindArticle 根据ID获取文章
// @Tags Article
// @Summary 根据ID获取文章
// @Description 根据文章ID获取文章详情
// @Produce application/json
// @Param id path integer true "文章ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /article/findArticle/{id} [get]
func (*ArticleApi) FindArticle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败", zap.Error(err))
		return response.FailWithMessage("获取id失败: "+err.Error(), c)
	}
	if articles, err := articleService.GetArticle(uint(id)); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "查询失败", c)
	} else {
		return response.OkWithData(articles, c)
	}
}

// GetArticleList 分页获取文章列表
// @Tags Article
// @Summary 分页获取文章列表
// @Description 分页获取文章列表，支持搜索和筛选
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
// @Router /article/getArticleList [get]
func (*ArticleApi) GetArticleList(c *fiber.Ctx) error {
	var pageInfo appReq.ArticleSearch
	_ = c.QueryParser(&pageInfo)
	IsImportant := c.QueryInt("is_important")
	pageInfo.IsImportant = IsImportant
	// log.Println("origin: ", c.Get("Origin"))
	if list, total, err := articleService.GetArticleInfoList(&pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "获取失败", c)

	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// PutArticleByIds 批量更新文章
// @Tags Article
// @Summary 批量更新文章
// @Description 批量更新文章的显示状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "文章ID列表"
// @Success 200 {object} response.Response{msg=string} "批量更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /article/putArticleByIds [put]
func (*ArticleApi) PutArticleByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	if err := c.BodyParser(&IDS); err != nil {
		// log.Println("ids 获取失败")
		global.LOG.Fatal("ids 获取失败", zap.Error(err))
		return response.FailWithMessage("ids 获取失败", c)
	}
	if err := articleService.PutArticleByIds(&IDS); err != nil {
		global.LOG.Error("批量更新失败!", zap.Error(err))
		return response.FailWithDetailed(map[string]string{
			"msg": err.Error(),
		}, "批量更新失败", c)

	} else {
		return response.OkWithMessage("批量更新成功", c)
	}
}

// GetArticleReading 获取文章阅读量
// @Tags Article
// @Summary 获取文章阅读量
// @Description 获取指定用户的文章阅读量统计
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{msg=string,data=object{reading_quantity=integer},code=integer} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Router /article/getArticleReading [get]
func (*ArticleApi) GetArticleReading(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	var id uint
	if userID != nil {
		id = userID.(uint)
	} else {
		id = 0
	}
	count, err := articleService.GetArticleReading(id)
	if err != nil {
		global.LOG.Error("获取阅读量失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{
			"msg": err.Error(),
		}, "获取阅读量失败", c)
	} else {
		return response.OkWithDetailed(fiber.Map{"reading_quantity": count}, "获取成功", c)
	}
}
