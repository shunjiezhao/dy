package handler

import (
	"context"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/pkg/redis"
	user3 "first/pkg/redis/user"
	"first/service/user/pack"
	user2 "first/service/user/service/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
)

func (s *UserServiceImpl) GetUsers(ctx context.Context, req *user.GetUserSRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return
	}
	klog.Infof("获取用户信息")

	resp.User, err = user3.GetUserInfo(redis.GetRedis(), ctx, req.Uuid, req.Id) // Redis 命中
	if err == nil {
		klog.Infof("redis hit")
		resp.Resp = pack.BuildBaseResp(errno.Success)
		return

	}

	if req.Uuid == 0 {
		resp.User, err = user2.NewGetUserService(ctx).GetUserS(req)
	} else {
		resp.User, err = user2.NewGetUserService(ctx).GetUserSWithLogin(req) // 登陆用户获取朋友列表, 关注字段
	}

	if err != nil {
		klog.Infof("获取用户信息失败%v", err)
		resp.Resp = pack.BuildBaseResp(errno.RemoteErr)
		return resp, nil
	}
	klog.Infof("操作成功")
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserResponse, err error) {
	log.Println("user rpc server: get user")
	resp = new(user.GetUserResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	users, err := user3.GetUserInfo(redis.GetRedis(), ctx, 0, []int64{req.Id}) // Redis 命中
	if err == nil && len(users) != 0 {
		klog.Infof("redis hit")
		resp.User = users[0]
		resp.Resp = pack.BuildBaseResp(errno.Success)
		return

	}

	resp.User, err = user2.NewGetUserService(ctx).GetUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	if len(req.UserName) == 0 || len(req.PassWord) == 0 {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return
	}
	resp = new(user.RegisterResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	resp.Id, err = user2.NewCreateUserService(ctx).CreateUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse,
	err error) {
	log.Println("user rpc server: check user")
	resp = &user.CheckUserResponse{} // 使用 new 里面的resp 不会初始化
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	resp.User, err = user2.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.AuthorizationFailedErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
