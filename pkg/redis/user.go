package redis

import (
	"context"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

const (
	// 用户的基本信息前缀
	userInfoPrefix         = "user."
	userNamePrefix         = userInfoPrefix + "name."          // string userid
	userFollowSetPrefix    = userInfoPrefix + "follow."        // set  userid
	userFollowerSetPrefix  = userInfoPrefix + "follower."      // set userid
	userLikeVideoPrefix    = userInfoPrefix + "like.video."    // zset 根据 排序 需要删除
	userPublishVideoPrefix = userInfoPrefix + "publish.video." // list
)

type UserKeyType int

const (
	Name UserKeyType = iota + 1
	Follow
	Follower
	LikeVideo
	PublishList
)

type UserRedis struct {
	ctx    context.Context
	userId int64
	redis  *redis.Client
}

func NewUserRedis(id int64, r *redis.Client, ctx context.Context) *UserRedis {
	return &UserRedis{userId: id, redis: r, ctx: ctx}
}

const (
	nameIdx          = iota //1. 用户名
	isFollowIdx             //2. 自己是否关注
	followCountIdx          //3. 关注数目
	followerCountIdx        //4. 粉丝数目
	idIdx                   //5. 用户id
)

// getVideoInfoLua
var getUserInfoLua = `
local ans = {}
for i, id in ipairs(ARGV) do
    local m = {}
    m[1]= redis.call('get', KEYS[1] .. id) -- 用户的姓名
    m[2]= redis.call('sismember',  KEYS[3] .. id, KEYS[4]) --是否关注 粉丝列表有没有自己 
    m[3]= redis.call('scard',  KEYS[2] .. id) -- 关注数目
    m[4]= redis.call('scard',  KEYS[3] .. id) -- 粉丝数目
	m[5] = tonumber(id) -- 用户id
    ans[i] = m
end
return ans
`

func (r *UserRedis) GetUserInfo(UserId []int64) (users []*handlers.User, err error) {
	args := make([]interface{}, len(users))
	for i := 0; i < len(UserId); i++ {
		args[i] = UserId[i]
	}
	result, err := r.redis.Eval(r.ctx, getFeedsInfoLua, []string{
		0: userNamePrefix,
		1: userFollowSetPrefix,
		2: userFollowerSetPrefix,
		3: intToString(r.userId),
	}, args...).Result()
	if err != nil {
		klog.Errorf("获取用户信息出错%v", err)
		return
	}
	return packUserHelper(result), nil
}
func packUserHelper(result interface{}) (users []*handlers.User) {
	defer func() {
		if a := recover(); a != nil {
			klog.Errorf("有 Key 过期", errno.RecordNotExistErr)
			users = nil
			return
		}
	}()

	res, ok := result.([]interface{})
	if !ok {
		return
	}
	for i := 0; i < len(res); i++ {
		val, ok := res[i].([]interface{})
		if !ok || len(val) != 5 {
			continue
		}
		user := &handlers.User{
			Id:            val[idIdx].(int64),
			Name:          val[nameIdx].(string),
			FollowCount:   val[followCountIdx].(int64),
			FollowerCount: val[followerCountIdx].(int64),
		}
		if val[followCountIdx].(int64) == 1 {
			user.IsFollow = true
		}
		users = append(users, user)
	}
	return users
}

var followUserListlua = `
local ans = {}
local ids = redis.call('smembers', KEYS[5]) -- 获取所有元素
for i, id in ipairs(ids) do
    local m = {}
    m[1]= redis.call('get', KEYS[1] .. id) -- 用户的姓名
    m[2]= redis.call('sismember',  KEYS[3] .. id, KEYS[4])  --是否关注 粉丝列表有没有自己 
    m[3]= redis.call('scard',  KEYS[2] .. id) -- 关注数目
    m[4]= redis.call('scard',  KEYS[3] .. id) -- 粉丝数目
	m[5] = tonumber(id) -- 用户id
    ans[i] = m
end
return ans
`

func (r *UserRedis) FollowUserList() (users []*handlers.User, err error) {
	return r.followHelper(Follow)
}
func (r *UserRedis) FollowerUserList() (users []*handlers.User, err error) {
	return r.followHelper(Follower)
}
func (r *UserRedis) followHelper(ty UserKeyType) (users []*handlers.User, err error) {
	eval := r.redis.Eval(r.ctx, followUserListlua, []string{
		0: userNamePrefix,
		1: userFollowSetPrefix,
		2: userFollowerSetPrefix,
		3: intToString(r.userId),
		4: GetKey(ty, r.userId),
	})
	result, err := eval.Result()
	if err != nil {
		return nil, err
	}
	return packUserHelper(result), nil

}

//LikeVideo 喜欢视频  更新视频信息点赞数, 将视频id push 到喜欢作品的list里面
func (r *UserRedis) LikeVideo(videoId int64) error {
	pipe := r.redis.Pipeline()
	// 加入视频的喜欢用户里面
	pipe.SAdd(r.ctx, GetVideoKey(LikeUser, videoId), r.userId)
	// 加入自己喜欢的视频列表
	pipe.ZAdd(r.ctx, GetKey(LikeVideo, r.userId), &redis.Z{
		Score:  float64(videoId),
		Member: videoId,
	}) // 点赞
	ans, err := pipe.Exec(r.ctx)
	klog.Infof("[喜欢视频操作]: 结果 %v", ans)
	return err
}

//FollowUser 关注操作
func (r *UserRedis) FollowUser(otherId int64) error {
	pipe := r.redis.Pipeline()
	// 1.增加follow粉丝列表
	pipe.SAdd(r.ctx, GetKey(Follower, otherId), r.userId)
	//	2.增加follower关注列表
	pipe.SAdd(r.ctx, GetKey(Follow, r.userId), otherId)
	ans, err := pipe.Exec(r.ctx)
	klog.Infof("[关注操作]: 结果 %v", ans)
	return err
}

//UnFollowUser 取消关注操作
func (r *UserRedis) UnFollowUser(otherId int64) error {
	pipe := r.redis.Pipeline()
	// 1.减少follow粉丝列表
	pipe.SRem(r.ctx, GetKey(Follower, otherId), r.userId)
	//	2.减少follower关注列表
	pipe.SRem(r.ctx, GetKey(Follow, r.userId), otherId)
	ans, err := pipe.Exec(r.ctx)
	klog.Infof("[关注操作]: 结果 %v", ans)
	return err
}

func GetKey(keyType UserKeyType, id int64) string {
	switch keyType {
	case Name:
		return userNameKey(id)
	case Follow:
		return UserFollowSetKey(id)
	case Follower:
		return userFollowerSetKey(id)
	case LikeVideo:
		return userLikeVideoKey(id)
	case PublishList:
		return userPublishVideoKey(id)
	}
	panic(errno.RedisKeyNotExistErr)
}

func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// 获取用户的Key

func userNameKey(id int64) string {
	return userNamePrefix + intToString(id)
}
func UserFollowSetKey(id int64) string {
	return userFollowSetPrefix + intToString(id)
}
func userFollowerSetKey(id int64) string {
	return userFollowerSetPrefix + intToString(id)
}

func userLikeVideoKey(id int64) string {
	return userLikeVideoPrefix + intToString(id)
}
func userPublishVideoKey(id int64) string {
	return userPublishVideoPrefix + intToString(id)
}
