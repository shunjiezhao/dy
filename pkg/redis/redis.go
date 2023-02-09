package redis

import (
	"first/pkg/constants"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

//init 初始Redis 连接
func init() {
	client = redis.NewClient(&redis.Options{
		Addr: constants.RedisDefaultURL,
	})
}

func GetRedis() *redis.Client {
	return client
}
