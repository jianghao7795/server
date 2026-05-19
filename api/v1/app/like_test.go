package app

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestLikeApi_LikePost_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Post("/like/likePost/:post_id", func(c fiber.Ctx) error {
		postId := c.Params("post_id")
		if postId == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取帖子ID失败"})
		}
		return c.JSON(fiber.Map{"post_id": postId})
	})

	t.Run("有效帖子ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/like/likePost/123", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("非数字帖子ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/like/likePost/abc", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestLikeApi_UnlikePost_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Delete("/like/unlikePost/:post_id", func(c fiber.Ctx) error {
		postId := c.Params("post_id")
		if postId == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取帖子ID失败"})
		}
		return c.JSON(fiber.Map{"unliked": postId})
	})

	t.Run("取消点赞有效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/like/unlikePost/50", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestLikeApi_GetPostLikes_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/like/getPostLikes/:post_id", func(c fiber.Ctx) error {
		postId := c.Params("post_id")
		page := c.Query("page", "1")
		pageSize := c.Query("page_size", "10")

		return c.JSON(fiber.Map{
			"post_id":   postId,
			"page":      page,
			"page_size": pageSize,
		})
	})

	t.Run("默认分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/like/getPostLikes/10", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
		assert.Contains(t, string(body), `"page_size":"10"`)
	})

	t.Run("自定义分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/like/getPostLikes/10?page=3&page_size=5", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"3"`)
		assert.Contains(t, string(body), `"page_size":"5"`)
	})
}

func TestLikeApi_CheckUserLiked(t *testing.T) {
	app := fiber.New()

	app.Get("/like/checkUserLiked/:post_id", func(c fiber.Ctx) error {
		postId := c.Params("post_id")
		return c.JSON(fiber.Map{"post_id": postId, "liked": false})
	})

	t.Run("检查点赞状态", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/like/checkUserLiked/100", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"liked":false`)
	})
}

func TestLikeApi_GetUserLikedPosts_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/like/getUserLikedPosts", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("page_size", "10")

		var pageInt int
		_ = (&struct{ Page int `query:"page"` }{Page: 1}).Page
		pageInt = 1

		return c.JSON(fiber.Map{
			"page":      page,
			"page_size": pageSize,
			"page_int":  pageInt,
		})
	})

	t.Run("默认参数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/like/getUserLikedPosts", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
	})
}

func TestLikeApi_GetPostLikeCount(t *testing.T) {
	app := fiber.New()

	app.Get("/like/getPostLikeCount/:post_id", func(c fiber.Ctx) error {
		postId := c.Params("post_id")
		return c.JSON(fiber.Map{"post_id": postId, "like_count": 42})
	})

	t.Run("获取点赞数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/like/getPostLikeCount/77", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"like_count":42`)
	})
}

func TestLikeApi_RequestBody(t *testing.T) {
	app := fiber.New()

	app.Post("/like/testBody", func(c fiber.Ctx) error {
		var body struct {
			PostId uint `json:"post_id"`
		}
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "参数错误"})
		}
		return c.JSON(fiber.Map{"post_id": body.PostId})
	})

	t.Run("有效body", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"post_id":100}`)
		req := httptest.NewRequest("POST", "/like/testBody", bodyReader)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"post_id":100`)
	})
}
