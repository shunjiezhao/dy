package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func GetFollowUserList(r *redis.Client, ctx context.Context, userId int64) (users []*userPb.User, err error) {
	return followHelper(r, ctx, Follow, userId)
}
func GetFollowerUserList(r *redis.Client, ctx context.Context, userId int64) (users []*userPb.User, err error) {
	return followHelper(r, ctx, Follower, userId)
}

var followUserListlua = `
if  redis.call('exists', KEYS[5]) == 0 then
	return
end
local ans = {}
local ids = redis.call('smembers', KEYS[5]) -- 获取所有元素
local ans = {}
local idx = 1
for i, id in ipairs(ids) do
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

func followHelper(r *redis.Client, ctx context.Context, ty UserKeyType, id int64) (users []*userPb.User,
	err error) {
	eval := r.Eval(ctx, followUserListlua, []string{
		0: UserNamePrefix,
		1: UserFollowSetPrefix,
		2: UserFollowerSetPrefix,
		3: strconv.FormatInt(id, 10),
		4: GetKey(ty, id),
	})
	result, err := eval.Result()
	if err != nil {
		return nil, err
	}
	return PackUserHelper(result)
}

func PackUserHelper(result interface{}) ([]*userPb.User, error) {
	res, ok := result.([]interface{})
	if !ok {
		klog.Errorf("断言出错")
		return nil, redis.Nil
	}
	users := make([]*userPb.User, 0, len(res))
	for i := 0; i < len(res); i++ {
		val := res[i].([]interface{})
		info := val[0].([]interface{})
		user := &userPb.User{
			Id:       val[2].(int64),
			UserName: info[0].(string),
		}
		if user.Id <= 0 {
			continue
		}

		if cnt, ok := info[1].(int64); ok {
			user.FollowCount = cnt
		} else {
			cnt, _ = strconv.ParseInt(info[1].(string), 10, 64)
			user.FollowCount = cnt
		}

		if cnt, ok := info[2].(int64); ok {
			user.FollowerCount = cnt
		} else {
			cnt, _ = strconv.ParseInt(info[2].(string), 10, 64)
			user.FollowerCount = cnt
		}
		if val[1].(int64) == 1 {
			user.IsFollow = true
		}
		users = append(users, user)
	}
	return users, nil
}

func WriteFollowList(r *redis.Client, ctx context.Context, me int64, users []*userPb.User) error {
	return writeFollowInfoHelper(r, ctx, me, users, Follow)
}
func WriteFollowerList(r *redis.Client, ctx context.Context, me int64, users []*userPb.User) error {
	return writeFollowInfoHelper(r, ctx, me, users, Follower)
}

// type 是将id加入到me的那个集合里面
func writeFollowInfoHelper(r *redis.Client, ctx context.Context, me int64, users []*userPb.User, ty UserKeyType) error {
	// 将用户信息写回
	pipe := r.Pipeline()
	for i := 0; i < len(users); i++ {
		r.SAdd(ctx, GetKey(ty, me), users[i].Id)
		hashKey := GetKey(UserInfo, users[i].Id)

		pipe.HMSet(ctx, hashKey, UserName, users[i].UserName, FollowCnt, users[i].FollowerCount,
			FollowerCnt, users[i].FollowCount)
		pipe.Expire(ctx, hashKey, ExpireTime)
	}
	result, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	klog.Infof("[WriteFollowerList]: %v", result)
	return nil
}

//DelUserFollowInfo form -> to 操作后删除各自的信息
func DelUserFollowInfo(r *redis.Client, ctx context.Context, from, to int64) error {
	pipe := r.Pipeline()
	pipe.Del(ctx, GetKey(UserInfo, from)) // from userinfo
	pipe.Del(ctx, GetKey(Follow, from))   // from follow set

	pipe.Del(ctx, GetKey(UserInfo, to)) // to
	pipe.Del(ctx, GetKey(Follower, to)) //
	result, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	klog.Infof("[DelUserFollowList]: %v", result)
	return err
}
