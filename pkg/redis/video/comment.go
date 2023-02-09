package video

import (
	"context"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"github.com/u2takey/go-utils/json"
)

// 返回评论列表
func GetCommentList(r *redis.Client, ctx context.Context, videoId int64, limit int64) ([]*handlers.Comment, error) {
	if limit > 10 {
		limit = 10
	}
	result, err := r.LRange(ctx, GetVideoKey(Comment, videoId), 0, limit).Result()
	klog.Infof("读取评论列表 %v", result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}
	ans := make([]*handlers.Comment, 0, len(result))
	for i := 0; i < len(result); i++ {
		comment := handlers.Comment{}
		err := json.Unmarshal([]byte(result[i]), &comment)
		if err != nil {
			klog.Errorf("解析失败 %v", err)
			r.Del(ctx, GetVideoKey(Comment, videoId)) // 删除
		}
		ans = append(ans, &comment)
	}
	return ans, nil
}

func WriteCommentList(r *redis.Client, ctx context.Context, videoId int64, comments []*handlers.Comment) error {
	comKey := GetVideoKey(Comment, videoId)
	pipe := r.Pipeline()
	for i := len(comments) - 1; i >= 0; i-- {
		pipe.LPush(ctx, comKey, comments[i]) // 左边是最新的 ,所以push
	}
	exec, err := pipe.Exec(ctx)
	klog.Infof("写入评论列表 %v", exec)
	if err != nil {
		return err
	}
	return nil
}
func DelCommetList(r *redis.Client, ctx context.Context, videoId int64) error {
	del := r.Del(ctx, GetVideoKey(Comment, videoId))
	klog.Infof("删除评论列表 %v", del)
	return del.Err()
}
