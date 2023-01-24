package main

import (
	"context"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/service"
	"log"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	log.Println("user rpc server")
	resp = new(user.RegisterResponse)
	resp.Id, err = service.NewCreateUserService(ctx).CreateUser(req)
	if err != nil {
		//TODO:检查是否为用户已存在
		resp.Resp.StatusCode = errno.UserAlreadyExistErr.ErrCode
		resp.Resp.StatusMsg = errno.UserAlreadyExistErr.ErrMsg
	}
	return
}
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse,
	err error) {
	log.Println("user rpc server: check user")
	resp = &user.CheckUserResponse{} // 不要用 new 否则里面的resp 不会初始化
	resp.Id, err = service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp.Resp.StatusCode = errno.AuthorizationFailedErr.ErrCode
		resp.Resp.StatusMsg = errno.AuthorizationFailedErr.ErrMsg
	}
	return
}
