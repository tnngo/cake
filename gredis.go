package cake

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/tnngo/cake/gredis"
	"go.uber.org/zap"
)

func loadRedis(rc *redisConfig) {
	gredis.RDB = redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
		PoolSize: rc.PoolSize,
	})

	if _, err := gredis.RDB.Ping(context.Background()).Result(); err != nil {
		zap.L().Error(err.Error())
		return
	} else {
		zap.L().Info("redis连接成功", zap.Reflect("config", rc))
	}

}
