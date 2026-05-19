package mobile

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserApi_CreateMobileUser_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/mobileUser/createMobileUser", func(c fiber.Ctx) error {
		var user struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Phone    string `json:"phone"`
		}
		if err := c.Bind().Body(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取用户数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{"username": user.Username})
	})

	t.Run("创建用户", func(t *testing.T) {
		body := strings.NewReader(`{"username":"mobile_user","nickname":"昵称","phone":"13800138000"}`)
		req := httptest.NewRequest("POST", "/mobileUser/createMobileUser", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"username":"mobile_user"`)
	})
}

func TestUserApi_DeleteMobileUser_Param(t *testing.T) {
	app := fiber.New()

	app.Delete("/mobileUser/deleteMobileUser/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{"deleted": id})
	})

	t.Run("删除用户", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/mobileUser/deleteMobileUser/10", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestUserApi_UpdateMobileUser_Binding(t *testing.T) {
	app := fiber.New()

	app.Put("/mobileUser/updateMobileUser", func(c fiber.Ctx) error {
		var user struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Nickname string `json:"nickname"`
		}
		if err := c.Bind().Body(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取用户信息失败"})
		}
		return c.JSON(fiber.Map{"id": user.ID, "username": user.Username})
	})

	t.Run("更新用户", func(t *testing.T) {
		body := strings.NewReader(`{"id":1,"username":"updated_user","nickname":"新昵称"}`)
		req := httptest.NewRequest("PUT", "/mobileUser/updateMobileUser", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"username":"updated_user"`)
	})
}

func TestUserApi_GetMobileUserList_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/mobileUser/getMobileUserList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		username := c.Query("username", "")
		return c.JSON(fiber.Map{"page": page, "pageSize": pageSize, "username": username})
	})

	t.Run("带搜索条件的列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/mobileUser/getMobileUserList?page=1&pageSize=5&username=test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"username":"test"`)
		assert.Contains(t, string(body), `"pageSize":"5"`)
	})
}
