package user

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"first/service/user/pack"
)

type CheckUserService struct {
	ctx context.Context
}

// NewCheckUserService new CreateNoteService
func NewCheckUserService(ctx context.Context) *CheckUserService {
	return &CheckUserService{ctx: ctx}
}

// CheckUser create note info
func (s *CheckUserService) CheckUser(req *user.CheckUserRequest) (*user.User, error) {
	println("rpc 响应开始调用")
	user, err := db.QueryUserByNamePwd(s.ctx, req.UserName, req.PassWord)
	if err != nil {
		return nil, err
	}

	return pack.User(user), err
}
