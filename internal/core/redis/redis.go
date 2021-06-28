package redis

import (
	"github.com/TheFootball/internal/configs"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func GetRedis() *redis.Client {
	if rdb != nil {
		return rdb
	}

	env := configs.GetEnv()

	rdb = redis.NewClient(&redis.Options{Addr: env.REDIS_ADDR, Password: env.REDIS_PW})
	return rdb
}
