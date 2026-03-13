package middleware

import (
	// "log"
	"server/model/common/response"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

type TokenBucket struct {
	capacity  int64      // 桶的容量
	rate      float64    // 令牌放入速率
	tokens    float64    // 当前令牌数量
	lastToken time.Time  // 上一次放令牌的时间
	mtx       sync.Mutex // 互斥锁
}

func (tb *TokenBucket) Allow() bool {
	tb.mtx.Lock()
	defer tb.mtx.Unlock()
	now := time.Now()
	// log.Println("now is sub oldnow ", now.Sub(tb.lastToken).Seconds())
	// 计算需要放的令牌数量
	tb.tokens = tb.tokens + tb.rate*now.Sub(tb.lastToken).Seconds()
	// log.Println("tb.tokens: ", tb.tokens)
	if max(tb.tokens, float64(tb.capacity)) != tb.tokens {
		tb.tokens = float64(tb.capacity)
	}

	// log.Println("tokens is {}", tb.tokens > 1)
	// 判断是否允许请求
	if tb.tokens >= 1 {
		tb.tokens--
		tb.lastToken = now
		return true
	} else {
		return false
	}
}

// LimitHandler 允许 连续访问次数
func LimitHandler(c fiber.Ctx) error {
	tb := &TokenBucket{
		capacity:  1000,
		rate:      10.0,
		tokens:    0,
		lastToken: time.Now(),
	}
	if !tb.Allow() {
		return response.FailWithDetailed(fiber.Map{"msg": "服务器需要休息一下，请等几分钟"}, "加载中", c)
	}
	return c.Next()
}
