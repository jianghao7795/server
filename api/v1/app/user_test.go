package app

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserApi_CreateUser_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/user/createUser", func(c fiber.Ctx) error {
		var user struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Phone    string `json:"phone"`
		}
		if err := c.Bind().Body(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{
			"username": user.Username,
			"nickname": user.Nickname,
		})
	})

	t.Run("创建用户", func(t *testing.T) {
		body := strings.NewReader(`{"username":"testuser","nickname":"测试","phone":"13800138000"}`)
		req := httptest.NewRequest("POST", "/user/createUser", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"username":"testuser"`)
	})
}

func TestUserApi_DeleteUser_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Delete("/user/deleteUser/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取id失败"})
		}
		return c.JSON(fiber.Map{"deleted": id})
	})

	t.Run("删除用户", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/user/deleteUser/123", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestUserApi_UpdateUser_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/user/updateUser/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		var user struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
		}
		if err := c.Bind().Body(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.JSON(fiber.Map{"id": id, "username": user.Username})
	})

	t.Run("更新用户", func(t *testing.T) {
		body := strings.NewReader(`{"username":"updated","nickname":"更新"}`)
		req := httptest.NewRequest("PUT", "/user/updateUser/10", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"username":"updated"`)
	})
}

func TestUserApi_GetUserList_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/user/getUserList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		return c.JSON(fiber.Map{"page": page, "pageSize": pageSize})
	})

	t.Run("默认分页", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user/getUserList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
		assert.Contains(t, string(body), `"pageSize":"10"`)
	})
}

func TestBaseMessageApi_CreateAndUpdate(t *testing.T) {
	app := fiber.New()

	app.Post("/base_message/create", func(c fiber.Ctx) error {
		var msg struct {
			Content string `json:"content"`
		}
		if err := c.Bind().Body(&msg); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{"created": msg.Content})
	})

	app.Put("/base_message/update/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		var msg struct {
			Content string `json:"content"`
		}
		_ = c.Bind().Body(&msg)
		return c.JSON(fiber.Map{"id": id, "content": msg.Content})
	})

	t.Run("创建基础消息", func(t *testing.T) {
		body := strings.NewReader(`{"content":"新消息"}`)
		req := httptest.NewRequest("POST", "/base_message/create", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"created":"新消息"`)
	})

	t.Run("更新基础消息", func(t *testing.T) {
		body := strings.NewReader(`{"content":"更新消息"}`)
		req := httptest.NewRequest("PUT", "/base_message/update/5", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(bodyBytes), `"id":"5"`)
		assert.Contains(t, string(bodyBytes), `"content":"更新消息"`)
	})
}

func TestFileUploadApi_RequestBody(t *testing.T) {
	app := fiber.New()

	app.Get("/file/uploadCheck", func(c fiber.Ctx) error {
		noSave := c.Query("noSave", "0")
		return c.JSON(fiber.Map{"noSave": noSave})
	})

	t.Run("noSave默认值", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/file/uploadCheck", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"noSave":"0"`)
	})

	t.Run("指定noSave", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/file/uploadCheck?noSave=1", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"noSave":"1"`)
	})
}
