package app

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestArticleApi_CreateArticle_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/article/createArticle", func(c fiber.Ctx) error {
		var article struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			UserId  int    `json:"user_id"`
		}
		if err := c.Bind().Body(&article); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{
			"title":   article.Title,
			"content": article.Content,
		})
	})

	t.Run("创建文章", func(t *testing.T) {
		body := strings.NewReader(`{"title":"测试标题","content":"测试内容","user_id":1}`)
		req := httptest.NewRequest("POST", "/article/createArticle", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"title":"测试标题"`)
	})
}

func TestArticleApi_DeleteArticle_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Delete("/article/deleteArticle/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" || id == "0" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取id失败"})
		}
		return c.JSON(fiber.Map{"id": id})
	})

	t.Run("删除有效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/article/deleteArticle/10", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("删除无效ID为0", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/article/deleteArticle/0", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestArticleApi_GetArticleList_QueryParams(t *testing.T) {
	app := fiber.New()

	app.Get("/article/getArticleList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		isImportant := c.Query("is_important", "0")
		title := c.Query("title", "")

		return c.JSON(fiber.Map{
			"page":         page,
			"pageSize":     pageSize,
			"is_important": isImportant,
			"title":        title,
		})
	})

	t.Run("带所有查询参数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/article/getArticleList?page=1&pageSize=20&is_important=1&title=go", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"is_important":"1"`)
		assert.Contains(t, string(body), `"title":"go"`)
	})

	t.Run("默认isImportant为0", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/article/getArticleList?page=2", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"is_important":"0"`)
	})
}

func TestArticleApi_UpdateArticle_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/article/updateArticle/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "参数错误"})
		}
		var article struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.Bind().Body(&article); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.JSON(fiber.Map{"id": id, "title": article.Title})
	})

	t.Run("更新文章", func(t *testing.T) {
		body := strings.NewReader(`{"title":"更新后的标题","content":"更新后的内容"}`)
		req := httptest.NewRequest("PUT", "/article/updateArticle/5", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"id":"5"`)
		assert.Contains(t, string(bodyBytes), `"title":"更新后的标题"`)
	})
}

func TestArticleApi_PutArticleByIds_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/article/putArticleByIds", func(c fiber.Ctx) error {
		var ids struct {
			Ids []int `json:"ids"`
		}
		if err := c.Bind().Body(&ids); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "ids获取失败"})
		}
		return c.JSON(fiber.Map{"count": len(ids.Ids)})
	})

	t.Run("批量更新", func(t *testing.T) {
		body := strings.NewReader(`{"ids":[10,20,30]}`)
		req := httptest.NewRequest("PUT", "/article/putArticleByIds", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"count":3`)
	})
}

func TestArticleApi_GetArticleReading(t *testing.T) {
	app := fiber.New()

	app.Get("/article/getArticleReading", func(c fiber.Ctx) error {
		userID := c.Locals("user_id")
		var id uint
		if userID != nil {
			id = userID.(uint)
		}
		return c.JSON(fiber.Map{"reading_quantity": 10, "user_id": id})
	})

	t.Run("有用户ID的请求", func(t *testing.T) {
		app2 := fiber.New()
		app2.Get("/test", func(c fiber.Ctx) error {
			c.Locals("user_id", uint(42))
			return c.JSON(fiber.Map{"reading_quantity": 5})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app2.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("无用户ID的请求默认0", func(t *testing.T) {
		app2 := fiber.New()
		app2.Get("/test2", func(c fiber.Ctx) error {
			userID := c.Locals("user_id")
			var id uint
			if userID != nil {
				id = userID.(uint)
			}
			return c.JSON(fiber.Map{"user_id": id})
		})

		req := httptest.NewRequest("GET", "/test2", nil)
		resp, _ := app2.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"user_id":0`)
	})
}
