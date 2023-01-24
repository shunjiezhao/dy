package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
)

type CheckUserService struct {
	ctx context.Context
}

// NewCheckUserService new CreateNoteService
func NewCheckUserService(ctx context.Context) *CheckUserService {
	return &CheckUserService{ctx: ctx}
}

// CheckUser create note info
func (s *CheckUserService) CheckUser(req *user.CheckUserRequest) (int64, error) {
	println("rpc 响应开始调用")
	user, err := db.QueryUserByName(s.ctx, req.UserName)
	if err != nil {
		return 0, err
	}

	return user.Uuid, err
}
