package example

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestCustomerApi_CreateCustomer_Binding(t *testing.T) {
	app := fiber.New()

	app.Post("/customer/customer", func(c fiber.Ctx) error {
		var customer struct {
			CustomerName       string `json:"customerName"`
			CustomerPhoneData  string `json:"customerPhoneData"`
		}
		if err := c.Bind().Body(&customer); err != nil {
			return c.Status(400).JSON(fiber.Map{"msg": "获取数据失败"})
		}
		return c.Status(200).JSON(fiber.Map{
			"customerName": customer.CustomerName,
			"phone":        customer.CustomerPhoneData,
		})
	})

	t.Run("创建客户", func(t *testing.T) {
		body := strings.NewReader(`{"customerName":"张三","customerPhoneData":"13800138000"}`)
		req := httptest.NewRequest("POST", "/customer/customer", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), `"customerName":"张三"`)
	})
}

func TestCustomerApi_DeleteCustomer_ParamValidation(t *testing.T) {
	app := fiber.New()

	app.Delete("/customer/customer/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" || id == "0" {
			return c.Status(400).JSON(fiber.Map{"msg": "id传递错误"})
		}
		return c.JSON(fiber.Map{"deleted": id})
	})

	t.Run("删除有效ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/customer/customer/10", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("删除ID为0", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/customer/customer/0", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestCustomerApi_GetCustomerList_Pagination(t *testing.T) {
	app := fiber.New()

	app.Get("/customer/customerList", func(c fiber.Ctx) error {
		page := c.Query("page", "1")
		pageSize := c.Query("pageSize", "10")
		return c.JSON(fiber.Map{"page": page, "pageSize": pageSize})
	})

	t.Run("分页查询", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/customer/customerList?page=2&pageSize=5", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), `"page":"2"`)
		assert.Contains(t, string(body), `"pageSize":"5"`)
	})
}

func TestCustomerApi_UpdateCustomer_Validation(t *testing.T) {
	app := fiber.New()

	app.Put("/customer/customer/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "0" {
			return c.Status(400).JSON(fiber.Map{"msg": "id不存在"})
		}
		var customer struct {
			ID           uint   `json:"id"`
			CustomerName string `json:"customerName"`
		}
		_ = c.Bind().Body(&customer)
		return c.JSON(fiber.Map{"updated": customer.CustomerName})
	})

	t.Run("更新客户", func(t *testing.T) {
		body := strings.NewReader(`{"id":5,"customerName":"更新的名称"}`)
		req := httptest.NewRequest("PUT", "/customer/customer/5", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
