package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	UserName     = "username"
	FollowCnt    = "followCount"
	FollowerCnt  = "followerCount"
	PublishCnt   = "publishCnt"
	LikeVideoCnt = "publishCnt"
	ExpireTime   = time.Minute
)

// getVideoInfoLua
var getUserInfoLua = `
local ans = {}
local idx = 1
for i, id in ipairs(ARGV) do
    local m = {}
	if  redis.call('exists', KEYS[1] .. id) == 1 then
		m[1]= redis.call('hmget', KEYS[1] .. id, KEYS[2], KEYS[3], KEYS[4]) -- 用户的姓名,  关注数目 ,  粉丝数目
		m[2]= redis.call('sismember',  KEYS[5] .. id, KEYS[6]) --是否关注 粉丝列表有没有自己 
		m[3] = tonumber(id) -- 用户id
    	ans[i] = m
	else 
		return -- 只要有一个不存在
	end
end
return ans
`

// mget
func GetUserInfo(r *redis.Client, ctx context.Context, me int64, userIds []int64) (users []*userPb.User, err error) {
	args := make([]interface{}, len(userIds))
	for i := 0; i < len(userIds); i++ {
		args[i] = userIds[i]
	}
	result, err := r.Eval(ctx, getUserInfoLua, []string{
		0: UserInfoPrefix,
		1: UserName,
		2: FollowCnt,
		3: FollowerCnt,
		4: UserFollowerSetPrefix,
		5: strconv.FormatInt(me, 10),
	}, args...).Result()
	if err != nil {
		klog.Errorf("获取用户信息出错%v", err)
		return
	}
	return PackUserHelper(result)
}
func WriteUserInfo(r *redis.Client, ctx context.Context, users []*userPb.User) (err error) {
	pipe := r.Pipeline()
	for i := 0; i < len(users); i++ {
		if users[i].Id > 0 {
			hashKey := GetKey(UserInfo, users[i].Id)
			pipe.HMSet(ctx, hashKey, UserName,
				users[i].UserName, FollowCnt, users[i].FollowCount,
				FollowerCnt, users[i].FollowerCount)
			pipe.Expire(ctx, hashKey, ExpireTime)
		}
	}
	res, err := pipe.Exec(ctx)
	if err != nil {
		klog.Errorf("[WriteUserInfo]: %v", err)
		return err
	}
	klog.Infof("[WriteUserInfo]: %v", res)
	return nil

}

func DelUserInfo(r *redis.Client, ctx context.Context, userId []int64) error {
	pipe := r.Pipeline()
	for i := 0; i < len(userId); i++ {
		pipe.Del(ctx, GetKey(UserInfo, userId[i]))
	}
	res, err := pipe.Exec(ctx)
	if err != nil {
		klog.Errorf("[DelUserInfo]: %v", err)
		return err
	}
	klog.Infof("[DelUserInfo]: %v", res)
	return nil
}
