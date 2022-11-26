package initialize

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"mybox/global"
)

func Redis() {
	redisCfg := global.BOX_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.BOX_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.BOX_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.BOX_REDIS = client
	}
}
