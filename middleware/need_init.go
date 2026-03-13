package middleware

import (
	"errors"
	"strings"

	global "server/model"

	"github.com/gofiber/fiber/v2"
)

// 处理跨域请求,支持options访问

func NeedInit(c *fiber.Ctx) error {
	var tables []string
	global.DB.Raw("show tables").Scan(&tables)
	if strings.Contains(strings.Join(tables, ""), "sys_users") {
		return c.Next()
	}
	return errors.New("没有初始化数据库")
	// 处理请求
}
