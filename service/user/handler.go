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
	resp = &user.CheckUserResponse{} // 使用 new 里面的resp 不会初始化
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
	resp = new(user.UserListResponse)
	resp.User, err = service.NewGetFollowerUserListService(ctx).GetFollowerUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, err
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.GetFollowListRequest) (resp *user.UserListResponse, err error) {
	// TODO: Your code here...
	resp = new(user.UserListResponse)
	resp.User, err = service.NewGetFollowUserListService(ctx).GetFollowUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, err
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// Follow implements the UserServiceImpl interface.
func (s *UserServiceImpl) Follow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse, err error) {
	log.Println("user rpc server: follow user")
	//TODO: 参数检查, to_user_id 是否合法
	//??? 如果再次关注会怎么样?
	resp = new(user.FollowResponse)
	_, err = service.NewFollowUserService(ctx).FollowUser(req)
	resp.Resp = pack.BuildBaseResp(err)

	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, err
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// UnFollow implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnFollow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse,
	err error) {
	log.Println("user rpc server: follow user")
	resp = new(user.FollowResponse)
	_, err = service.NewUnFollowUserService(ctx).UnFollowUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, err
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
