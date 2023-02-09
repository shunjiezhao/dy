package video

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/kitex_gen/video"
	"first/pkg/redis/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"time"
)

// WriteFavVideo 将 用户的喜欢列表中的视频信息更新
func WriteFavVideo(r *redis.Client, ctx context.Context, videos []*video.Video, userId int64) error {
	ids := make([]interface{}, 0, len(videos))

	pipe := r.Pipeline()
	for i := 0; i < len(videos); i++ {
		writeVideoInfoHelper(pipe, ctx, videos)
		ids = append(ids, videos[i].Id)
	}

	pipe.SAdd(ctx, user.GetKey(user.LikeVideo, userId), ids...).Result()
	exec, err := pipe.Exec(ctx)
	klog.Infof("WriteFavVideo: %v", exec)
	if err != nil {
		return err
	}

	return nil
}

// getVideoInfoLua
var getFavInfoLua = `
if redis.call('exists', KEYS[1]) == 0 then 
	return
end
local ans = {}

local videoIds = redis.call('smembers', KEYS[1]) -- 用户 喜欢视频集合
for i, id in ipairs(videoIds) do -- 遍历
	local m = {}
    if redis.call('exists', KEYS[2] .. id) == 1 then -- 如果不存在当前视频信息 说明信息过期, 从新获取
		m[1] = redis.call('hmget', KEYS[2] .. id, KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7], KEYS[8])
		m[2] = tonumber(id)
		ans[i] = m
	else
		return -- 有一个信息不存在
    end
end
return ans
`

func GetFavVideoList(r *redis.Client, ctx context.Context, userId int64) ([]*video.Video, error) {
	result, err := r.Eval(ctx, getFavInfoLua, []string{
		0: user.GetKey(user.LikeVideo, userId), // 用户的喜欢列表
		1: videoInfoPrefix,                     // key 剑指亲追
		// 2
		2: videoPlayUrlKey, // url key
		3: videoCoverUrKey,
		4: videoAuthorKey,
		5: videoCommentCntKey,
		6: videoFavCntKey,
		7: videoTitleKey,
	}).Result()
	klog.Infof("[GetFavVideoList]: %+v", result)
	if err != nil {
		return nil, err
	}
	res, ok := result.([]interface{})
	if !ok {
		return nil, redis.Nil
	}
	vides := make([]*video.Video, 0, len(res))
	for i := 0; i < len(res); i++ {
		val, ok := res[i].([]interface{})
		if !ok || len(val) != 2 {
			klog.Errorf("断言失败 %v", res[i])
			return nil, redis.Nil
		}

		info, ok := val[0].([]interface{})
		if !ok || len(info) != 5 {
			klog.Errorf("断言失败 %v", val)
			return nil, redis.Nil
		}
		vides = append(vides, &video.Video{
			Id:            val[1].(int64),
			Author:        &userPb.User{Id: info[0].(int64)},
			PlayUrl:       info[1].(string),
			CoverUrl:      info[2].(string),
			FavoriteCount: info[3].(int64),
			CommentCount:  info[4].(int64),
			IsFavorite:    true,
			Title:         info[5].(string),
		})
	}
	return vides, err
}
func writeVideoInfoHelper(pipe redis.Pipeliner, ctx context.Context, v []*video.Video) {

	for i := 0; i < len(v); i++ {
		key := GetVideoHashKey(v[i].Id)
		if v[i].Author.Id > 0 && v[i].Id > 0 &&
			v[i].PlayUrl != "" && v[i].CoverUrl != "" {
			pipe.HMSet(ctx, key.VideoHashKey(),
				key.CoverUrlKey(), v[i].CoverUrl,
				key.PlayUrlKey(), v[i].PlayUrl,
				key.FavCountKey(), v[i].FavoriteCount,
				key.Author(), v[i].Author.Id,
				key.CommentCount(), v[i].CommentCount,
				key.Title(), v[i].Title,
			)

		}

		pipe.Expire(ctx, key.VideoHashKey(), time.Second)
	}
	klog.Infof("更新视频信息操作")
}

func DelFavVideoList(r *redis.Client, ctx context.Context, userId int64) error {
	_, err := r.Del(ctx, user.GetKey(user.LikeVideo, userId)).Result()
	klog.Infof("[DelFavVideoList]: %v", err)
	return err
}
