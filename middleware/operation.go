// Package middleware 提供中间件功能，包括操作记录、认证等
package middleware

import (
	"encoding/json"
	global "server-fiber/model"
	"server-fiber/model/system"
	"server-fiber/utils"
	"strconv"
	"strings"
	"time"

	// json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	systemService "server-fiber/service/system"

	"go.uber.org/zap"
)

// operationRecordService 操作记录服务实例
var operationRecordService = new(systemService.OperationRecordService)

// OperationRecord 操作记录中间件
// 用于记录用户的API操作历史，包括请求信息、响应状态、执行时间等
// 支持GET请求的查询参数记录和POST/PUT等请求的请求体记录
// 自动识别请求来源（后端管理、前端API、移动端）并记录相应的端口类型
func OperationRecord(c *fiber.Ctx) error {
	var body []byte
	var userId int
	// var err error

	// 根据请求方法处理请求体数据
	if c.Method() != fiber.MethodGet {
		// 非GET请求直接获取请求体
		body = c.Request().Body()
	} else {
		// GET请求解析查询参数并转换为JSON格式
		query := c.OriginalURL()
		split := strings.Split(query, "?")
		if len(split) > 1 {
			// 解析查询参数
			splitI := strings.Split(split[1], "&")
			m := make(map[string]string)
			for _, v := range splitI {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			// 将查询参数转换为JSON格式存储
			body, _ = json.Marshal(&m)
		}

	}
	// 获取JWT claims信息
	claims, err := utils.GetClaims(c)
	if err != nil {
		return c.Status(403).JSON(map[string]string{"msg": err.Error()})
	}
	// fmt.Printf("%v\n", claims)

	// 获取用户ID，优先从JWT claims中获取，否则从请求头中获取
	if claims.BaseClaims.ID != 0 {
		userId = int(claims.BaseClaims.ID)
	} else {
		// 从请求头中获取用户ID
		id, err := strconv.Atoi(c.Get("x-user-id"))
		if err != nil {
			userId = 0
		}
		userId = id
	}
	// 获取请求路径并判断请求来源类型
	pathURL := c.Path()
	isBackend := system.Backend
	switch {
	case strings.HasPrefix(pathURL, "/backend"):
		// 后端管理接口
		isBackend = system.Backend
	case strings.HasPrefix(pathURL, "/api"):
		// 前端API接口
		isBackend = system.Frontend
	case strings.HasPrefix(pathURL, "/mobile"):
		// 移动端接口
		isBackend = system.Mobile
	default:
		// 默认为后端管理接口
		isBackend = system.Backend
	}

	// 创建操作记录对象
	record := system.SysOperationRecord{
		Ip:       c.IP(),              // 客户端IP地址
		Method:   c.Method(),          // HTTP请求方法
		Path:     pathURL,             // 请求路径
		Agent:    c.Get("User-Agent"), // 用户代理信息
		Body:     string(body),        // 请求体内容
		UserID:   userId,              // 用户ID
		TypePort: isBackend,           // 请求来源类型
	}

	// 处理文件上传请求，对请求体进行长度限制
	if strings.Contains(c.Get("Content-Type"), "multipart/form-data") {
		if len(record.Body) > 512 {
			record.Body = "File or Length out of limit"
		}
	}
	// 使用defer确保在请求处理完成后记录操作信息
	defer func() {
		// 记录响应状态码
		record.Status = c.Response().StatusCode()

		// 如果是服务器内部错误，记录错误信息
		if record.Status == fiber.StatusInternalServerError {
			record.ErrorMessage = string(c.Response().Body())
		} else {
			record.ErrorMessage = ""
		}

		// 计算请求处理延迟时间
		record.Latency = time.Since(time.Now())

		// 记录响应内容
		record.Resp = string(c.Response().Body())

		// 保存操作记录到数据库
		if err := operationRecordService.CreateSysOperationRecord(&record); err != nil {
			global.LOG.Error("create operation record error:", zap.Error(err))
		}
	}()

	// 继续执行下一个中间件或处理器
	return c.Next()
}
