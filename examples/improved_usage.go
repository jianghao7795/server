package examples

import (
	"server-fiber/middleware"
	"server-fiber/model/app"
	"server-fiber/model/common/request"
	"server-fiber/model/common/response"
	"server-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

// ExampleAPI demonstrates improved API implementation
type ExampleAPI struct {
	articleService *utils.CRUDBase[app.Article]
}

// NewExampleAPI creates a new example API
func NewExampleAPI() *ExampleAPI {
	return &ExampleAPI{
		articleService: utils.NewCRUDBase[app.Article](nil), // Pass actual DB instance
	}
}

// CreateArticleExample demonstrates improved article creation
func (api *ExampleAPI) CreateArticleExample(c *fiber.Ctx) error {
	var article app.Article

	// Use validation middleware for cleaner code
	if err := middleware.ValidationMiddlewareInstance.ValidateBody(c, &article); err != nil {
		return err
	}

	// Use CRUD base for database operations
	if err := api.articleService.Create(&article); err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "创建文章", err)
	}

	return response.OkWithId("创建成功", article.ID, c)
}

// GetArticleExample demonstrates improved article retrieval
func (api *ExampleAPI) GetArticleExample(c *fiber.Ctx) error {
	// Use validation middleware for ID validation
	id, err := middleware.ValidationMiddlewareInstance.ValidateID(c)
	if err != nil {
		return err
	}

	// Use CRUD base for database operations
	article, err := api.articleService.GetByID(id)
	if err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "获取文章", err)
	}

	return response.OkWithData(article, c)
}

// UpdateArticleExample demonstrates improved article update
func (api *ExampleAPI) UpdateArticleExample(c *fiber.Ctx) error {
	var article app.Article

	// Validate request body
	if err := middleware.ValidationMiddlewareInstance.ValidateBody(c, &article); err != nil {
		return err
	}

	// Validate ID parameter
	id, err := middleware.ValidationMiddlewareInstance.ValidateID(c)
	if err != nil {
		return err
	}

	article.ID = id

	// Use CRUD base for database operations
	if err := api.articleService.Update(&article); err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "更新文章", err)
	}

	return response.OkWithMessage("更新成功", c)
}

// DeleteArticleExample demonstrates improved article deletion
func (api *ExampleAPI) DeleteArticleExample(c *fiber.Ctx) error {
	// Validate ID parameter
	id, err := middleware.ValidationMiddlewareInstance.ValidateID(c)
	if err != nil {
		return err
	}

	// Use CRUD base for database operations
	if err := api.articleService.Delete(id); err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "删除文章", err)
	}

	return response.OkWithMessage("删除成功", c)
}

// GetArticleListExample demonstrates improved article listing
func (api *ExampleAPI) GetArticleListExample(c *fiber.Ctx) error {
	var pageInfo request.PageInfo

	// Parse query parameters
	if err := c.QueryParser(&pageInfo); err != nil {
		return utils.ErrorHandlerInstance.HandleValidationError(c, "分页参数", err)
	}

	// Set default values
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}

	// Use CRUD base for database operations
	list, total, err := api.articleService.GetList(pageInfo)
	if err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "获取文章列表", err)
	}

	return response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// DeleteArticlesExample demonstrates batch deletion
func (api *ExampleAPI) DeleteArticlesExample(c *fiber.Ctx) error {
	var ids request.IdsReq

	// Validate request body
	if err := middleware.ValidationMiddlewareInstance.ValidateBody(c, &ids); err != nil {
		return err
	}

	// Use CRUD base for database operations
	if err := api.articleService.DeleteByIDs(ids); err != nil {
		return utils.ErrorHandlerInstance.HandleAPIError(c, "批量删除文章", err)
	}

	return response.OkWithMessage("批量删除成功", c)
}
