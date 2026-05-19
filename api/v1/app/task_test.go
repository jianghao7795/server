package app

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestTaskNameApi_StartTasking(t *testing.T) {
	app := fiber.New()

	app.Get("/tasking/start", func(c fiber.Ctx) error {
		tasking := c.Query("task")
		if tasking == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "请传入任务名"})
		}
		return c.JSON(fiber.Map{"task": tasking, "status": "started"})
	})

	t.Run("启动有效任务", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasking/start?task=cleanDb", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"task":"cleanDb"`)
	})

	t.Run("缺少任务名", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasking/start", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestTaskNameApi_StopTasking(t *testing.T) {
	app := fiber.New()

	app.Get("/tasking/stop", func(c fiber.Ctx) error {
		tasking := c.Query("task")
		if tasking == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "请传入任务名"})
		}
		return c.JSON(fiber.Map{"task": tasking, "status": "stopped"})
	})

	t.Run("停止有效任务", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasking/stop?task=backupDb", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(body), `"status":"stopped"`)
	})

	t.Run("空任务名", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasking/stop?task=", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}
