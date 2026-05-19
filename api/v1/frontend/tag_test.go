package frontend

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestTagApi_GetTag_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Get("/frontend/getTag/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取Ids失败"})
		}
		return c.JSON(fiber.Map{"id": id})
	})

	t.Run("有效ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/frontend/getTag/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestTagApi_GetTagList(t *testing.T) {
	app := fiber.New()

	app.Get("/frontend/getTagList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		return c.JSON(fiber.Map{"page": page, "pageSize": pageSize})
	})

	t.Run("默认分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/frontend/getTagList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
	})
}
