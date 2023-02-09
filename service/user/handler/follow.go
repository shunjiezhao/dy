package handler

import (
	"context"
	"first/kitex_gen/user"
	"first/pkg/errno"
	"first/pkg/redis"
	user3 "first/pkg/redis/user"
	"first/service/user/pack"
	"first/service/user/service/follow"
	user2 "first/service/user/service/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
	"time"
)

// 处理流程
// 1. 参数
// 2. redis
// 3. db

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.GetFollowerListRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	resp.User, err = user3.GetFollowerUserList(redis.GetRedis(), ctx, req.Id)
	if err == nil || len(resp.User) != 0 {
		goto redisHit
	}
	resp.User, err = follow.NewGetFollowerUserListService(ctx).GetFollowerUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return resp, nil
	}
	err = user3.WriteFollowerList(redis.GetRedis(), ctx, req.Id, resp.User)
	if err != nil {
		klog.Errorf("写入失败")
		go user3.WriteFollowerList(redis.GetRedis(), ctx, req.Id, resp.User) // 在写一次
	}
redisHit:
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.GetFollowListRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	resp.User, err = user3.GetFollowUserList(redis.GetRedis(), ctx, req.Id)
	if err == nil || len(resp.User) != 0 {
		goto redisHit
	}

	resp.User, err = follow.NewGetFollowUserListService(ctx).GetFollowUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return resp, nil
	}
	err = user3.WriteFollowList(redis.GetRedis(), ctx, req.Id, resp.User) // 写入 redis
	if err != nil {
		klog.Errorf("写入失败")
		go user3.WriteFollowList(redis.GetRedis(), ctx, req.Id, resp.User) // 在写一次
	}
redisHit:
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// 写操作处理流程
//1. 延时双删
//2. 更新 DB

// Follow implements the UserServiceImpl interface.
func (s *UserServiceImpl) Follow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse, err error) {

	log.Println("user rpc server: follow user")
	//??? 如果再次关注会怎么样?
	resp = new(user.FollowResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	go func() {
		err = user3.DelUserFollowInfo(redis.GetRedis(), ctx, req.FromUserId, req.ToUserId)
		if err != nil {
			klog.Errorf("redis 删除失败")

		}
		time.Sleep(500 * time.Millisecond)
		err = user3.DelUserFollowInfo(redis.GetRedis(), ctx, req.FromUserId, req.ToUserId)
		if err != nil {
			klog.Errorf("redis 删除失败")

		}
	}()

	_, err = user2.NewFollowUserService(ctx).FollowUser(req)
	resp.Resp = pack.BuildBaseResp(err)

	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// UnFollow implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnFollow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse,
	err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		err = user3.DelUserFollowInfo(redis.GetRedis(), ctx, req.FromUserId, req.ToUserId)
		if err != nil {
			klog.Errorf("redis 删除失败")

		}

	}()
	log.Println("user rpc server: follow user")
	resp = new(user.FollowResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	err = user3.DelUserFollowInfo(redis.GetRedis(), ctx, req.FromUserId, req.ToUserId)
	if err != nil {
		klog.Errorf("redis 删除失败")

	}
	_, err = user2.NewUnFollowUserService(ctx).UnFollowUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
