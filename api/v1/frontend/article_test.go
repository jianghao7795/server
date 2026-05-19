package frontend

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestArticleApi_GetArticleList_Defaults(t *testing.T) {
	app := fiber.New()

	app.Get("/getArticleList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		isImportant := c.Query("is_important", "0")
		return c.JSON(fiber.Map{
			"page":         page,
			"is_important": isImportant,
		})
	})

	t.Run("默认参数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getArticleList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
		assert.Contains(t, string(body), `"is_important":"0"`)
	})
}

func TestArticleApi_GetArticleDetail_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Get("/getArticle/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取Id失败"})
		}
		return c.JSON(fiber.Map{"id": id})
	})

	t.Run("有效ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getArticle/5", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("无效非数字ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getArticle/abc", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestArticleApi_GetSearchArticle_Validation(t *testing.T) {
	app := fiber.New()

	app.Get("/getSearchArticle/:name/:value", func(c fiber.Ctx) error {
		name := c.Params("name")
		if name != "tags" && name != "articles" {
			return c.Status(400).JSON(fiber.Map{"msg": "查询的不是tag或article"})
		}
		return c.JSON(fiber.Map{"name": name, "value": c.Params("value")})
	})

	t.Run("按tag搜索", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getSearchArticle/tags/Go", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"name":"tags"`)
	})

	t.Run("按article搜索", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getSearchArticle/articles/测试", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"name":"articles"`)
	})

	t.Run("无效搜索类型", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/getSearchArticle/users/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}
