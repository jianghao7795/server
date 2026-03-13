package middleware

import (
	"context"
	"errors"
	global "server/model"
	"server/model/common/response"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type LimitConfig struct {
	// GenerationKey 根据业务生成key 下面CheckOrMark查询生成
	GenerationKey func(c *fiber.Ctx) string
	// 检查函数,用户可修改具体逻辑,更加灵活
	CheckOrMark func(key string, expire int, limit int) error
	// Expire key 过期时间
	Expire int
	// Limit 周期时间
	Limit int
}

func (l LimitConfig) LimitWithTime(c *fiber.Ctx) error {
	if err := l.CheckOrMark(l.GenerationKey(c), l.Expire, l.Limit); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": response.ERROR, "msg": err})
	}
	return c.Next()
}

// DefaultGenerationKey 默认生成key
func DefaultGenerationKey(c *fiber.Ctx) string {
	return "Limit" + c.IP()
}

func DefaultCheckOrMark(key string, expire int, limit int) (err error) {
	// 判断是否开启redis
	if global.REDIS == nil {
		return err
	}
	if err = SetLimitWithTime(key, limit, time.Duration(expire)*time.Second); err != nil {
		global.LOG.Error("limit", zap.Error(err))
	}
	return err
}

func DefaultLimit(c *fiber.Ctx) error {
	return LimitConfig{
		GenerationKey: DefaultGenerationKey,
		CheckOrMark:   DefaultCheckOrMark,
		Expire:        global.CONFIG.System.LimitTimeIP,
		Limit:         global.CONFIG.System.LimitCountIP,
	}.LimitWithTime(c)
}

// SetLimitWithTime 设置访问次数
func SetLimitWithTime(key string, limit int, expiration time.Duration) error {
	count, err := global.REDIS.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		pipe := global.REDIS.TxPipeline()
		pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, expiration)
		_, err = pipe.Exec(context.Background())
		return err
	} else {
		// 次数
		if times, err := global.REDIS.Get(context.Background(), key).Int(); err != nil {
			return err
		} else {
			if times >= limit {
				if t, err := global.REDIS.PTTL(context.Background(), key).Result(); err != nil {
					return errors.New("请求太过频繁，请稍后再试")
				} else {
					return errors.New("请求太过频繁, 请 " + t.String() + " 秒后尝试")
				}
			} else {
				return global.REDIS.Incr(context.Background(), key).Err()
			}
		}
	}
}
