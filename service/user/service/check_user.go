package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
)

type CheckUserService struct {
	ctx context.Context
}

// NewCreateNoteService new CreateNoteService
func NewCheckUserService(ctx context.Context) *CheckUserService {
	return &CheckUserService{ctx: ctx}
}

// CreateUser create note info
func (s *CheckUserService) ChekcUser(req *user.CheckUserRequest) (int64, error) {
	println("rpc 响应开始调用")
	user, err := db.QueryUser(s.ctx, req.UserName)
	//TODO: error 是否加工
	if err != nil {
		return 0, err
	}
	return user.Uuid, err
}
