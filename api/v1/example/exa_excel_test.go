package example

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestExcelApi_ExportExcel_Validation(t *testing.T) {
	app := fiber.New()

	app.Post("/excel/exportExcel", func(c fiber.Ctx) error {
		var info struct {
			FileName string `json:"fileName"`
		}
		if err := c.Bind().Body(&info); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "参数错误"})
		}
		if info.FileName == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "请传参数filename"})
		}
		if strings.Contains(info.FileName, "..") {
			return c.Status(400).JSON(fiber.Map{"msg": "包含非法字符"})
		}
		return c.Status(200).JSON(fiber.Map{"fileName": info.FileName})
	})

	t.Run("有效文件名", func(t *testing.T) {
		body := strings.NewReader(`{"fileName":"export.xlsx"}`)
		req := httptest.NewRequest("POST", "/excel/exportExcel", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"fileName":"export.xlsx"`)
	})

	t.Run("空文件名", func(t *testing.T) {
		body := strings.NewReader(`{"fileName":""}`)
		req := httptest.NewRequest("POST", "/excel/exportExcel", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("包含非法字符..", func(t *testing.T) {
		body := strings.NewReader(`{"fileName":"../etc/passwd"}`)
		req := httptest.NewRequest("POST", "/excel/exportExcel", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestExcelApi_LoadExcel(t *testing.T) {
	app := fiber.New()

	app.Get("/excel/loadExcel", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"total": 10, "page": 1, "pageSize": 999})
	})

	req := httptest.NewRequest("GET", "/excel/loadExcel", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestExcelApi_DownloadTemplate_Validation(t *testing.T) {
	app := fiber.New()

	app.Get("/excel/downloadTemplate", func(c fiber.Ctx) error {
		fileName := c.Query("fileName")
		if fileName == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "请传入fileName"})
		}
		return c.JSON(fiber.Map{"fileName": fileName})
	})

	t.Run("有效模板名", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/excel/downloadTemplate?fileName=template.xlsx", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("缺少文件名", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/excel/downloadTemplate", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestExcelApi_GetFileList(t *testing.T) {
	app := fiber.New()

	app.Get("/excel/getFileList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		return c.JSON(fiber.Map{"page": page, "pageSize": pageSize})
	})

	t.Run("文件列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/excel/getFileList", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"1"`)
	})
}
