package middleware

import (
	"server/model/common/response"
	service "server/service/system"
	"server/utils"

	"github.com/gofiber/fiber/v3"
)

var casbinService = new(service.CasbinService)

// 拦截器
func CasbinHandler(c fiber.Ctx) error {
	waitUse, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("token 错误", c)
	}
	// 获取请求的PATH
	obj := c.Path()
	// 获取请求方法
	act := c.Method()
	// 获取用户的角色
	sub := waitUse.AuthorityId
	e := casbinService.Casbin()
	// 判断策略中是否存在
	success, err := e.Enforce(sub, obj, act)
	// log.Println("error is ", err, success, obj, act, sub)
	// if global.CONFIG.System.Env == "develop" || success {
	if err != nil {
		return response.FailWithMessage403("验证失败", c)
	}
	if success {
		return c.Next()
	} else {
		// 上传文件 由于是ajxs 必须返回403 错误 才能展示错误信息
		return response.FailWithMessage403("权限不足", c)
	}
}
