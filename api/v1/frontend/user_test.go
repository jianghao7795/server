package frontend

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestUser_RegisterUser_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/frontend/register", func(c fiber.Ctx) error {
		var user struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			RePassword  string `json:"rePassword"`
		}
		if err := c.Bind().Body(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		if user.Password != user.RePassword {
			return c.Status(400).JSON(fiber.Map{"msg": "密码不一致"})
		}
		return c.Status(200).JSON(fiber.Map{"username": user.Username, "msg": "注册成功"})
	})

	t.Run("注册成功", func(t *testing.T) {
		body := strings.NewReader(`{"username":"newuser","password":"123456","rePassword":"123456"}`)
		req := httptest.NewRequest("POST", "/frontend/register", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"msg":"注册成功"`)
	})

	t.Run("密码不一致", func(t *testing.T) {
		body := strings.NewReader(`{"username":"newuser","password":"123456","rePassword":"654321"}`)
		req := httptest.NewRequest("POST", "/frontend/register", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestUser_UpdatePassword_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/frontend/updatePassword", func(c fiber.Ctx) error {
		var data struct {
			NewPassword       string `json:"newPassword"`
			RepeatNewPassword string `json:"repeatNewPassword"`
		}
		if err := c.Bind().Body(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		if data.NewPassword != data.RepeatNewPassword {
			return c.Status(400).JSON(fiber.Map{"msg": "密码不一致"})
		}
		return c.Status(200).JSON(fiber.Map{"msg": "重置密码成功"})
	})

	t.Run("密码一致", func(t *testing.T) {
		body := strings.NewReader(`{"newPassword":"newpass","repeatNewPassword":"newpass"}`)
		req := httptest.NewRequest("PUT", "/frontend/updatePassword", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("密码不一致", func(t *testing.T) {
		body := strings.NewReader(`{"newPassword":"newpass","repeatNewPassword":"different"}`)
		req := httptest.NewRequest("PUT", "/frontend/updatePassword", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}
