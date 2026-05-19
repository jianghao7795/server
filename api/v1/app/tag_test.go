package app

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestTagApi_CreateTag_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/tag/createTag", func(c fiber.Ctx) error {
		var tag struct {
			Name   string `json:"name"`
			Status int    `json:"status"`
		}
		if err := c.Bind().Body(&tag); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{"name": tag.Name, "status": tag.Status})
	})

	t.Run("创建标签", func(t *testing.T) {
		body := strings.NewReader(`{"name":"Go语言","status":1}`)
		req := httptest.NewRequest("POST", "/tag/createTag", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"name":"Go语言"`)
		assert.Contains(t, string(bodyBytes), `"status":1`)
	})
}

func TestTagApi_GetTagList_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/tag/getTagList", func(c fiber.Ctx) error {
		p := c.Query("page", "1")
		ps := c.Query("pageSize", "10")
		return c.JSON(fiber.Map{
			"page":     p,
			"pageSize": ps,
		})
	})

	t.Run("默认分页参数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tag/getTagList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"page":"1"`)
	})

	t.Run("指定分页参数", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tag/getTagList?page=2&pageSize=5", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"2"`)
		assert.Contains(t, string(body), `"pageSize":"5"`)
	})
}

func TestTagApi_DeleteTag_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Delete("/tag/deleteTag/:id", func(c fiber.Ctx) error {
		idParam := c.Params("id")
		if idParam == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "id不能为空"})
		}
		return c.JSON(fiber.Map{"deleted": idParam})
	})

	t.Run("删除有效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tag/deleteTag/42", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"deleted":"42"`)
	})
}

func TestTagApi_UpdateTag_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/tag/updateTag", func(c fiber.Ctx) error {
		var tag struct {
			ID     uint   `json:"id"`
			Name   string `json:"name"`
			Status int    `json:"status"`
		}
		if err := c.Bind().Body(&tag); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.JSON(fiber.Map{
			"id":     tag.ID,
			"name":   tag.Name,
			"status": tag.Status,
		})
	})

	t.Run("更新标签", func(t *testing.T) {
		body := strings.NewReader(`{"id":1,"name":"更新后的标签","status":0}`)
		req := httptest.NewRequest("PUT", "/tag/updateTag", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"name":"更新后的标签"`)
	})
}
