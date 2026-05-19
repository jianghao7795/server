package mobile

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestLoginApi_Login_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/mobile/login", func(c fiber.Ctx) error {
		var l struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Realname string `json:"realname"`
		}
		if err := c.Bind().Body(&l); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取登录数据失败"})
		}
		if l.Username == "" || l.Password == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "用户名或密码不能为空"})
		}
		return c.Status(200).JSON(fiber.Map{"username": l.Username})
	})

	t.Run("登录成功", func(t *testing.T) {
		body := strings.NewReader(`{"username":"testuser","password":"123456","realname":"张"}`)
		req := httptest.NewRequest("POST", "/mobile/login", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"username":"testuser"`)
	})

	t.Run("缺少用户名", func(t *testing.T) {
		body := strings.NewReader(`{"username":"","password":"123456"}`)
		req := httptest.NewRequest("POST", "/mobile/login", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestLoginApi_GetUserInfo(t *testing.T) {
	app := fiber.New()

	app.Get("/mobile/getUserInfo", func(c fiber.Ctx) error {
		userId := c.Get("user_id", "")
		if userId == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "获取失败"})
		}
		return c.JSON(fiber.Map{"userId": userId})
	})

	t.Run("有user_id请求头", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/mobile/getUserInfo", nil)
		req.Header.Set("user_id", "123")
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"userId":"123"`)
	})

	t.Run("无user_id请求头", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/mobile/getUserInfo", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestLoginApi_UpdateMobileUser(t *testing.T) {
	app := fiber.New()

	app.Put("/mobile/updateMobileUser", func(c fiber.Ctx) error {
		userId := c.Get("user_id", "")
		if userId == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "更新失败"})
		}
		var data struct {
			Field string `json:"field"`
			Value string `json:"value"`
		}
		_ = c.Bind().Body(&data)
		return c.JSON(fiber.Map{"userId": userId, "field": data.Field, "value": data.Value})
	})

	t.Run("更新用户信息", func(t *testing.T) {
		body := strings.NewReader(`{"field":"nickname","value":"新昵称"}`)
		req := httptest.NewRequest("PUT", "/mobile/updateMobileUser", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("user_id", "42")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"field":"nickname"`)
		assert.Contains(t, string(bodyBytes), `"userId":"42"`)
	})

	t.Run("无user_id", func(t *testing.T) {
		body := strings.NewReader(`{"field":"nickname","value":"新昵称"}`)
		req := httptest.NewRequest("PUT", "/mobile/updateMobileUser", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestLoginApi_UpdatePassword(t *testing.T) {
	app := fiber.New()

	app.Put("/mobile/updatePassword", func(c fiber.Ctx) error {
		var data struct {
			ID          uint   `json:"id"`
			Password    string `json:"password"`
			NewPassword string `json:"newPassword"`
		}
		if err := c.Bind().Body(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.JSON(fiber.Map{"id": data.ID, "newPassword": data.NewPassword})
	})

	t.Run("更新密码", func(t *testing.T) {
		body := strings.NewReader(`{"id":1,"password":"old","newPassword":"new"}`)
		req := httptest.NewRequest("PUT", "/mobile/updatePassword", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"id":1`)
	})
}
