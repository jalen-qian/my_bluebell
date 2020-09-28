package redis

import (
	"github.com/go-redis/redis"
	"github.com/web_app_base/settings"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err = rdb.Ping().Result()
	return
}
