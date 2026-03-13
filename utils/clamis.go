package utils

import (
	"errors"
	global "server/model"
	systemReq "server/model/system/request"
	"strings"

	"github.com/gofiber/fiber/v3"
	uuid "github.com/google/uuid"
)

// GetClaims 从fiber的Context中获取JWT token并解析出CustomClaims结构体
func GetClaims(c fiber.Ctx) (*systemReq.CustomClaims, error) {
	tokenString := c.Get("Authorization", "") // 从请求头中获取token字符串
	// fmt.Println("tokenString: ", tokenString)
	if tokenString == "" {
		return nil, errors.New("没有登录，请重新登录")
	}
	token := strings.Replace(tokenString, "Bearer ", "", 1) // 去除token字符串中的"Bearer "前缀
	if token == "" {
		return nil, errors.New("token 失效， 请重新登录") // 若token为空则返回401错误响应
	}
	j := NewJWT()                      // 初始化JWT对象
	claims, err := j.ParseToken(token) // 解析token获取CustomClaims结构体和可能的错误
	if err != nil {
		global.LOG.Error("从fiber的Context中获取从jwt解析信息失败, 请检查请求头是否存在token且claims是否为规定结构") // 若解析出错则记录错误日志
	}
	return claims, err // 返回解析出的CustomClaims结构体和错误
}

// GetUserID 从fiber的Context中获取从jwt解析出来的用户ID
func GetUserID(c fiber.Ctx) (uint, error) {
	waitUse, ok := c.Locals("claims").(*systemReq.CustomClaims)
	if !ok {
		// 处理类型断言失败的情况，例如返回错误信息或默认值
		return 0, nil
	}

	if waitUse.BaseClaims.ID == 0 {
		// 处理ID为0的情况，可能需要调用GetClaims函数重新获取claims
		if cl, err := GetClaims(c); err != nil {
			// 处理获取claims失败的情况，例如返回错误信息或默认值
			return 0, err
		} else {
			return uint(cl.BaseClaims.ID), nil
		}
	} else {
		return waitUse.BaseClaims.ID, nil
	}
}

// GetUserUuid 从fiber的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c fiber.Ctx) uuid.UUID {
	claims := c.Locals("claims")
	waitUse, ok := claims.(*systemReq.CustomClaims)
	if !ok {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	}
	if waitUse.BaseClaims.ID == 0 {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		return waitUse.UUID
	}
}

// 从fiber的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c fiber.Ctx) (string, error) {
	claims := c.Locals("claims")
	waitUse, ok := claims.(*systemReq.CustomClaims)
	if !ok {
		if cl, err := GetClaims(c); err != nil {
			return "", err
		} else {
			return cl.AuthorityId, nil
		}
	}
	if waitUse.BaseClaims.ID == 0 {
		if cl, err := GetClaims(c); err != nil {
			return "", err
		} else {
			return cl.AuthorityId, nil
		}
	} else {
		return waitUse.AuthorityId, nil
	}
}

// 从fiber的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c fiber.Ctx) *systemReq.CustomClaims {
	waitUse, ok := c.Locals("claims").(*systemReq.CustomClaims)
	if !ok {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	}
	if waitUse.BaseClaims.ID == 0 {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	}
	return waitUse
}
