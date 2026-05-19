package mobile

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestRegisterMobile_Register_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/mobile/register", func(c fiber.Ctx) error {
		var data struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.Bind().Body(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{"username": data.Username, "msg": "注册成功"})
	})

	t.Run("注册用户", func(t *testing.T) {
		body := strings.NewReader(`{"username":"newmobile","password":"password123"}`)
		req := httptest.NewRequest("POST", "/mobile/register", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), "注册成功")
	})
}
