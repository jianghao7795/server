package system

import (
	"server/model/common/response"
	"server/model/system"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type JwtApi struct{}

// @Tags Jwt
// @Summary jwt加入黑名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "jwt加入黑名单"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /jwt/jsonInBlacklist [post]
func (j *JwtApi) JsonInBlacklist(c fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	token := strings.Replace(tokenString, "Bearer ", "", 1)
	if token == "" {
		return response.FailWithMessage401("token 失效， 请重新登录", 3, nil, c)
	}
	jwt := system.JwtBlacklist{Jwt: token}
	if err := jwtService.JsonInBlacklist(jwt); err != nil {
		return response.FailWithMessage("jwt作废失败", 3, err, c)
	} else {
		return response.OkWithMessage("jwt作废成功", c)
	}
}
