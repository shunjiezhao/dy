package redis

import (
	"context"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
)

const (
	// 视频的基本信息前缀
	videoInfoPrefix = "video."
	// map play_url, cover_url, author_id
	videoAuthorPrefix   = videoInfoPrefix + "author."
	videoPlayUrlPrefix  = videoInfoPrefix + "play.url."
	videoCoverUrlPrefix = videoInfoPrefix + "cover.url."

	videoCommentPrefix  = videoInfoPrefix + "comment."   // list
	videoLikeUserPrefix = videoInfoPrefix + "like.user." //  set

	videoFeedListKey = videoInfoPrefix + "feed.list"
)

type VideoKeyType int

const (
	Author VideoKeyType = iota + 1
	PlayUrl
	CoverUrl
	Comment
	LikeUser
	FeedList
)

type videoRedis struct {
	ctx     context.Context
	videoId int64
	redis   *redis.Client
}

//PublishVideo 发布视频
func PublishVideo(r *redis.Client, ctx context.Context, videoId int64) error {
	err := r.LPush(ctx, GetVideoKey(FeedList, 0), videoId).Err()
	klog.Infof("[发布视频]: ", err)
	return err
}

const (
	playUrlIdx = iota
	coverUrlIdx
	commentIdx
	favCountIdx
	isFavouriteIdx
	videoIdIdx
	authorIdx
)

// getVideoInfoLua
var getFeedsInfoLua = `
local ans = {}
local videoIds = redis.call('lrange', KEYS[1], ARGV[1], ARGV[2])
for i, id in ipairs(videoIds) do
    local m = {}
    m[1]= redis.call('get', KEYS[2] .. id)
    m[2]= redis.call('get',  KEYS[3] .. id)
    m[3]= redis.call('llen',  KEYS[4] .. id)
    m[4]= redis.call('scard',  KEYS[5] .. id)
    m[5] = redis.call('sismember',  KEYS[5] .. id, ARGV[3])
	m[6] = tonumber(id)
	m[7] = tonumber(redis.call('get', KEYS[6] .. id))
    ans[i] = m
end
return ans
`

//GetFeedsVideo 获取视频流
func GetFeedsVideo(r *redis.Client, ctx context.Context, limit int, userId int64) ([]*handlers.Video, error) {
	if limit > 10 {
		limit = 10
	}
	var (
		videos []*handlers.Video
	)
	// 1.获取视频 url
	// 2 封面url
	// 3.获取点赞数
	// 4. 评论数
	val, err := r.Eval(ctx, getVideoInfoLua, []string{
		0: videoFeedListKey,
		1: videoPlayUrlPrefix,
		2: videoCoverUrlPrefix,
		3: videoCommentPrefix,
		4: videoLikeUserPrefix,
		5: videoAuthorPrefix,
	}, 0, limit, userId).Result()

	res, ok := val.([]interface{})
	if err != nil || !ok { // 找不到 去获取
		klog.Infof("redis 获取信息失败")
		return nil, err
	}
	if len(res) == 0 {
		return nil, errno.RecordNotExistErr
	}

	for _, tval := range res {
		val, ok := tval.([]interface{})
		if !ok || len(val) < 6 {
			continue
		}
		v := &handlers.Video{
			Author:        &handlers.User{Id: val[authorIdx].(int64)}, // 视频信息
			PlayUrl:       val[playUrlIdx].(string),
			CoverUrl:      val[coverUrlIdx].(string),
			CommentCount:  val[commentIdx].(int64),
			FavoriteCount: val[favCountIdx].(int64),
		}
		if val[isFavouriteIdx].(int64) == 1 {
			v.IsFavorite = true
		}
		v.Id = val[videoIdIdx].(int64)
		videos = append(videos, v)
	}
	return videos, nil
}

// getVideoInfoLua 左到右依次是 play_url
var getVideoInfoLua = `
local m = {}
m[1]= redis.call('get', KEYS[1])
m[2]= redis.call('get', KEYS[2])
m[3]= redis.call('llen', KEYS[3])
m[4]= redis.call('scard', KEYS[4])
m[5]= redis.call('sismember' , KEYS[4], ARGV[1])
m[6]= tonumber(redis.call('get', KEYS[5]))
return  m
`

//GetVideoInfo 获取视频信息
func GetVideoInfo(r *redis.Client, ctx context.Context, videoId int64, userId int64) *handlers.Video {
	// 如果没有该视频的信息就直接返回
	var (
		video handlers.Video
	)
	// 1.获取视频 url
	// 2 封面url
	// 3.获取点赞数
	// 4. 评论数
	result, err := r.Eval(ctx, getVideoInfoLua, []string{
		0: GetVideoKey(PlayUrl, videoId),
		1: GetVideoKey(CoverUrl, videoId),
		2: GetVideoKey(Comment, videoId),
		3: GetVideoKey(LikeUser, videoId),
		4: GetVideoKey(Author, videoId),
	}, userId).Result()
	if err != nil {
		if err != redis.Nil {
			klog.Infof("获取视频信息出错")
			return nil
		}
	}
	res, ok := result.([]interface{})
	if !ok || len(res) != 6 {
		klog.Infof("获取视频信息出错")
		return nil
	}

	video.PlayUrl = res[0].(string)
	video.CoverUrl = res[1].(string)
	video.CommentCount = res[2].(int64)
	video.FavoriteCount = res[3].(int64)
	video.Author = &handlers.User{Id: res[5].(int64)}

	if res[4].(int64) == 1 {
		video.IsFavorite = true
	}

	video.Id = videoId
	return &video
}

func GetVideoKey(keyType VideoKeyType, videoId int64) string {
	switch keyType {
	case Author:
		return videoAuthorKey(videoId)
	case PlayUrl:
		return videoPlayUrlKey(videoId)
	case CoverUrl:
		return videoCoverUrlKey(videoId)
	case Comment:
		return videoCommentKey(videoId)
	case LikeUser:
		return videoLikeUserKey(videoId)
	case FeedList:
		return videoFeedListKey
	}
	panic(errno.RedisKeyNotExistErr)
}

// 获取用户的Key

func videoAuthorKey(id int64) string {
	return videoAuthorPrefix + intToString(id)
}
func videoPlayUrlKey(id int64) string {
	return videoPlayUrlPrefix + intToString(id)
}
func videoCoverUrlKey(id int64) string {
	return videoCoverUrlPrefix + intToString(id)
}

func videoCommentKey(id int64) string {
	return videoCommentPrefix + intToString(id)
}
func videoLikeUserKey(id int64) string {
	return videoLikeUserPrefix + intToString(id)
}
