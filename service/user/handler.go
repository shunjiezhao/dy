package main

import (
	"context"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/pack"
	"first/service/user/service"
	"log"
)

//TODO:
// 1.对于 string 类型 进行 SQL 注入检查
// 2.参数检查
// 3.检查是否该返回用户已存在的错误

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserResponse, err error) {
	log.Println("user rpc server: get user")
	//TODO: 参数检查
	resp = new(user.GetUserResponse)
	resp.User, err = service.NewGetUserService(ctx).GetUser(req)
	if err != nil {
		//TODO:检查是否为用户已存在
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	log.Println("user rpc server")
	//TODO: 参数检查
	if len(req.UserName) == 0 || len(req.PassWord) == 0 {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return
	}
	resp = new(user.RegisterResponse)
	resp.Id, err = service.NewCreateUserService(ctx).CreateUser(req)
	if err != nil {
		//TODO:检查是否为用户已存在
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse,
	err error) {
	log.Println("user rpc server: check user")
	//TODO: 参数检查
	resp = &user.CheckUserResponse{} // 不要用 new 否则里面的resp 不会初始化
	resp.Id, err = service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.AuthorizationFailedErr)
		return
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.GetFollowerListRequest) (resp *user.UserListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.GetFollowListRequest) (resp *user.UserListResponse, err error) {
	// TODO: Your code here...
	return
}
