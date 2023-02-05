package user

import (
	"context"
	"first/kitex_gen/user"
	"first/pkg/util"
	"first/service/user/model/db"
)

type CreateUserService struct {
	ctx context.Context
}

// NewCreateUserService new CreateNoteService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// CreateUser create note info
func (s *CreateUserService) CreateUser(req *user.RegisterRequest) (int64, error) {
	println("rpc 响应开始调用")
	userModel := &db.User{
		Uuid:     util.NextVal(),
		UserName: req.UserName,
		Password: encryptPassWord(req.PassWord),
		NickName: req.UserName,
	}
	return db.CreateUser(s.ctx, userModel)
}

//TODO: 加密密码
func encryptPassWord(passWord string) string {
	return passWord
}
