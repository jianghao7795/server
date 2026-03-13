package init_load

import (
	"context"
	global "server/model"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Redis() error {
	// 初始化时可以指定连接redis的读写超时时间，默认都是3s
	redisCfg := global.CONFIG.Redis
	s := []string{redisCfg.Addr, ":", redisCfg.Port}
	client := redis.NewClient(&redis.Options{
		Addr:            strings.Join(s, ""), // redis服务ip:port
		Password:        redisCfg.Password,   // redis的认证密码
		DB:              redisCfg.DB,         // 连接的database
		ConnMaxIdleTime: 30 * time.Minute,    // 默认Idle超时时间
		PoolSize:        100,                 // 连接池
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	// redisExample2(client, ctx)

	if err != nil {
		global.LOG.Error("redis connect ping failed, err:", zap.Error(err))
		return err
	} else {
		global.LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.REDIS = client
	}
	return err
}

// func RedisExample2(rdb *redis.Client, ctx context.Context) {
// 	zsetKey := "language_rank"
// 	languages := []redis.Z{
// 		{Score: 90.0, Member: "Golang"},
// 		{Score: 98.0, Member: "Java"},
// 		{Score: 95.0, Member: "Python"},
// 		{Score: 97.0, Member: "JavaScript"},
// 		{Score: 99.0, Member: "C/C++"},
// 	}
// 	// ZADD
// 	num, err := rdb.ZAdd(ctx, zsetKey, languages...).Result()
// 	if err != nil {
// 		log.Printf("zadd failed, err:%v\n", err)
// 		return
// 	}
// 	log.Printf("zadd %d succ.\n", num)

// 	// 把Golang的分数加10
// 	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
// 	if err != nil {
// 		log.Printf("zincrby failed, err:%v\n", err)
// 		return
// 	}
// 	log.Printf("Golang's score is %f now.\n", newScore)

// 	// 取分数最高的3个
// 	ret, err := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Result()
// 	if err != nil {
// 		fmt.Printf("zrevrange failed, err:%v\n", err)
// 		return
// 	}
// 	for _, z := range ret {
// 		fmt.Println(z.Member, z.Score)
// 	}

// 	// 取95~100分的
// 	op := redis.ZRangeBy{
// 		Min: "99",
// 		Max: "100",
// 	}
// 	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, &op).Result()
// 	if err != nil {
// 		fmt.Printf("zrangebyscore failed, err:%v\n", err)
// 		return
// 	}
// 	for _, z := range ret {
// 		fmt.Println(z.Member, z.Score)
// 	}

// 	vals, err := rdb.Keys(ctx, "Golang").Result()

// 	if err != nil {
// 		log.Printf("keys failed, err:%v\n", err)
// 		return
// 	}
// 	log.Println(vals)
// }
