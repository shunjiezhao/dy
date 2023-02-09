package video

import (
	"first/pkg/errno"
	"strconv"
)

// video 基本信息
// title
// comments count
// like user count
// author id

const (
	// 视频的基本信息前缀
	videoInfoPrefix = "video."
	// map play_url, cover_url, author_id
	videoPlayUrlKey = "playurl"
	videoCoverUrKey = "coverurl"
	videoAuthorKey  = "author"
	videoTitleKey   = "title"

	videoCommentCntKey = "commentCount"
	videoFavCntKey     = "favCount"

	videoCommentPrefix = "video.comment." // 视频评论列表
)

type VideoKeyType int

const (
	VideoInfo VideoKeyType = iota + 1
	Comment
)

// getVideoInfoLua
var getFeedsInfoLua = `
if redis.call('exists', KEYS[1]) == 0 then 
	return
end

local ans = {}
local videoIds = redis.call('lrange', KEYS[1], ARGV[1], ARGV[2])
local idx = 1
for i, id in ipairs(videoIds) do
    local m = {}
    if redis.call('exists', KEYS[2] .. id) == 1 then
		m[1]= redis.call('get', KEYS[2] .. id)
		m[2]= redis.call('get',  KEYS[3] .. id)
		m[3]= redis.call('llen',  KEYS[4] .. id)
		m[4]= redis.call('scard',  KEYS[5] .. id)
		m[5] = redis.call('sismember',  KEYS[5] .. id, ARGV[3])
		m[6] = tonumber(id)
		m[7] = tonumber(redis.call('get', KEYS[6] .. id)) -- 作者
	    ans[idx] = m
        idx = idx + 1
    end
end
return ans
`

type VideoHash struct {
	videoId string
}

func GetVideoHashKey(videoId int64) VideoHash {
	return VideoHash{
		intToString(videoId),
	}
}
func (h VideoHash) VideoHashKey() string {
	return videoInfoPrefix + h.videoId
}
func (h VideoHash) Author() string {
	return videoAuthorKey
}
func (h VideoHash) CommentCount() string {
	return videoCommentCntKey
}
func (h VideoHash) FavCountKey() string {
	return videoFavCntKey
}
func (h VideoHash) PlayUrlKey() string {
	return videoPlayUrlKey
}
func (h VideoHash) CoverUrlKey() string {
	return videoCoverUrKey
}
func (h VideoHash) Title() string {
	return videoCoverUrKey
}

func GetVideoKey(keyType VideoKeyType, videoId int64) string {
	switch keyType {
	case Comment:
		return videoCommentKey(videoId)
	case VideoInfo:
		return videoHashInfoKey(videoId)
	}
	panic(errno.RedisKeyNotExistErr)
}

// 获取用户的Key

func videoCommentKey(id int64) string {
	return videoCommentPrefix + intToString(id)
}
func videoHashInfoKey(id int64) string {
	return videoInfoPrefix + intToString(id)
}

func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
