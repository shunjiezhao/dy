package redus

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

const RedisDefaultURL = "127.0.0.1:6379"

type Metadata struct { //元数据结构体 包含名称、size、哈希值
	Name string //对象的名字
	Size int64  // 大小 根据content-Length
	Hash string // 利用 SHA-256 hash 得来的 ”“ 代表没有
}

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: RedisDefaultURL,
	})
	if client == nil {
		panic("redis 初始化失败")
	}
	//TODO: scan file

}

// JudgeIsExist 判断hash 值是否存在
func JudgeIsExist(ctx context.Context, hash string) bool {
	result, err := client.Get(ctx, hash).Result()
	if err != nil {
		if redis.Nil != err {
			klog.Errorf("redis 获取出错")
			return false
		}
	}

	i, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		klog.Errorf("解析数字出错, 检查redis value设置错误, key 为: %v", hash)
		return false
	}
	return i > 0
}

func SetKey(ctx context.Context, hash string, size int64) error {
	pipe := client.Pipeline()
	pipe.Set(ctx, hash, size, 0)
	if res, err := pipe.Exec(ctx); err != nil {
		klog.Info(res)
		return err
	}

	return nil
}

func GetKey(ctx context.Context, hash string) int64 {
	res, err := client.Get(ctx, hash).Result()
	klog.Info(res)
	if err != nil {
		return 0
	}
	size, _ := strconv.ParseInt(res, 10, 64)
	return size
}
