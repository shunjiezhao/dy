package redis

import (
	"first/pkg/constants"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	redis *redis.Client
}

//InitRedis 初始Redis 连接
func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: constants.RedisDefaultURL,
	})
	return client
}

func NewReids() *Redis {
	return &Redis{redis: InitRedis()}
}

func IsRedisError(err error) bool {
	return err == redis.Nil
}
