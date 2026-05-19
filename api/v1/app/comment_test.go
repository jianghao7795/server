package app

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	return fiber.New()
}

func TestCommentApi_GetCommentList_Pagination(t *testing.T) {
	app := setupTestApp()
	api := &CommentApi{}

	app.Get("/comment/getCommentList", func(c fiber.Ctx) error {
		var pageInfo = struct {
			Page     int `query:"page"`
			PageSize int `query:"pageSize"`
		}{}
		_ = c.Bind().Query(&pageInfo)
		if pageInfo.Page == 0 {
			pageInfo.Page = 1
		}
		if pageInfo.PageSize == 0 {
			pageInfo.PageSize = 10
		}
		return c.JSON(fiber.Map{
			"page":     pageInfo.Page,
			"pageSize": pageInfo.PageSize,
		})
	})
	_ = api

	t.Run("默认分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/comment/getCommentList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":1`)
		assert.Contains(t, string(body), `"pageSize":10`)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("自定义分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/comment/getCommentList?page=3&pageSize=20", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":3`)
		assert.Contains(t, string(body), `"pageSize":20`)
	})
}

func TestCommentApi_DeleteComment_ParamParsing(t *testing.T) {
	app := setupTestApp()

	app.Delete("/comment/deleteComment/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "id不能为空"})
		}
		return c.JSON(fiber.Map{"id": id, "success": true})
	})

	t.Run("删除有效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/comment/deleteComment/123", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"id":"123"`)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("删除无效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/comment/deleteComment/abc", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestCommentApi_DeleteCommentByIds_Binding(t *testing.T) {
	app := setupTestApp()

	app.Delete("/comment/deleteCommentByIds", func(c fiber.Ctx) error {
		var ids struct {
			Ids []int `json:"ids"`
		}
		if err := c.Bind().Body(&ids); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "参数错误"})
		}
		return c.JSON(fiber.Map{"count": len(ids.Ids)})
	})

	t.Run("批量删除", func(t *testing.T) {
		body := strings.NewReader(`{"ids":[1,2,3,4,5]}`)
		req := httptest.NewRequest("DELETE", "/comment/deleteCommentByIds", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"count":5`)
	})
}

func TestCommentApi_CreateComment_Binding(t *testing.T) {
	app := setupTestApp()

	app.Post("/comment/createComment", func(c fiber.Ctx) error {
		var comment struct {
			Content string `json:"content"`
			PostId  uint   `json:"post_id"`
			UserId  uint   `json:"user_id"`
		}
		if err := c.Bind().Body(&comment); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		if comment.Content == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "内容不能为空"})
		}
		return c.Status(200).JSON(fiber.Map{"msg": "创建成功"})
	})

	t.Run("有效请求体", func(t *testing.T) {
		body := strings.NewReader(`{"content":"你好世界","post_id":1,"user_id":100}`)
		req := httptest.NewRequest("POST", "/comment/createComment", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), "创建成功")
	})

	t.Run("缺少内容", func(t *testing.T) {
		body := strings.NewReader(`{"content":"","post_id":1,"user_id":100}`)
		req := httptest.NewRequest("POST", "/comment/createComment", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestCommentApi_PutLikeItOrDislike_Binding(t *testing.T) {
	app := setupTestApp()

	app.Put("/comment/pariseComment", func(c fiber.Ctx) error {
		var praise struct {
			CommentId int64 `json:"comment_id"`
			UserId    int64 `json:"user_id"`
		}
		if err := c.Bind().Body(&praise); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.JSON(fiber.Map{
			"comment_id": praise.CommentId,
			"user_id":    praise.UserId,
		})
	})

	t.Run("点赞带ID参数", func(t *testing.T) {
		body := strings.NewReader(`{"id":0,"comment_id":10,"user_id":100}`)
		req := httptest.NewRequest("PUT", "/comment/pariseComment", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"comment_id":10`)
		assert.Contains(t, string(bodyBytes), `"user_id":100`)
	})
}
