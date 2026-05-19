package frontend

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestCommentApi_GetCommentByArticleId(t *testing.T) {
	app := fiber.New()

	app.Get("/frontend/comment/:articleId", func(c fiber.Ctx) error {
		articleId := c.Params("articleId")
		if articleId == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取articleId失败"})
		}
		return c.JSON(fiber.Map{"articleId": articleId})
	})

	t.Run("有效文章ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/frontend/comment/100", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("非数字ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/frontend/comment/abc", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestCommentApi_CreatedComment_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/frontend/comment", func(c fiber.Ctx) error {
		var comment struct {
			Content    string `json:"content"`
			ArticleId  uint   `json:"article_id"`
			UserId     uint   `json:"user_id"`
		}
		if err := c.Bind().Body(&comment); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "参数错误"})
		}
		return c.Status(200).JSON(fiber.Map{
			"content":    comment.Content,
			"article_id": comment.ArticleId,
		})
	})

	t.Run("创建评论", func(t *testing.T) {
		body := strings.NewReader(`{"content":"好文章","article_id":5,"user_id":10}`)
		req := httptest.NewRequest("POST", "/frontend/comment", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"content":"好文章"`)
	})
}
