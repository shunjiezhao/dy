package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"first/service/user/pack"
)

type GetUserService struct {
	ctx context.Context
}

// NewGetUserService new CreateNoteService
func NewGetUserService(ctx context.Context) *GetUserService {
	return &GetUserService{ctx: ctx}
}

// GetUser create note info
func (s *GetUserService) GetUser(req *user.GetUserRequest) (*user.User, error) {
	println("rpc 响应开始调用")
	user, err := db.QueryUserById(s.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return pack.User(user), err
}
