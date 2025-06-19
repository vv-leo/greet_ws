package initialize

import (
	"git.exclouds.org/ww/greet_ws/global"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func RedisInitialize() {
	redisCfg := global.SERVER_CONFIG.RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default CONTACT_DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.LOGGER.Error("redis connect ping failed, err:", zap.Any("err", err))
	} else {
		global.LOGGER.Info("redis connect ping resp:", zap.String("pong", pong))
		global.REDIS = client
	}
}
