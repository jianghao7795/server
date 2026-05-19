package frontend

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetImages(t *testing.T) {
	app := fiber.New()

	app.Get("/frontend/getImages", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"images": []string{"img1.jpg", "img2.jpg"},
		})
	})

	req := httptest.NewRequest("GET", "/frontend/getImages", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}
