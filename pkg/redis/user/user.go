package user

import (
	"context"
	"first/pkg/errno"
	"github.com/go-redis/redis/v8"
	"strconv"
)

const (
	// 用户的基本信息前缀
	InfoPrefix         = "user."
	NamePrefix         = InfoPrefix + "name."          // string userid
	FollowSetPrefix    = InfoPrefix + "follow."        // set  userid
	FollowerSetPrefix  = InfoPrefix + "follower."      // set userid
	LikeVideoPrefix    = InfoPrefix + "like.video."    // zset 根据 排序 需要删除
	PublishVideoPrefix = InfoPrefix + "publish.video." // list
)

type UserKeyType int

const (
	UserInfo UserKeyType = iota + 1
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

func GetKey(keyType UserKeyType, id int64) string {
	switch keyType {
	case UserInfo:
		return userInfoKey(id)
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

// 获取用户的Key

func userInfoKey(id int64) string {
	return InfoPrefix + intToString(id)
}
func UserFollowSetKey(id int64) string {
	return FollowSetPrefix + intToString(id)
}
func userFollowerSetKey(id int64) string {
	return FollowerSetPrefix + intToString(id)
}

func userLikeVideoKey(id int64) string {
	return LikeVideoPrefix + intToString(id)
}
func userPublishVideoKey(id int64) string {
	return PublishVideoPrefix + intToString(id)
}
func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
