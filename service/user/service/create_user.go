package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
)

type CreateUserService struct {
	ctx context.Context
}

// NewCreateNoteService new CreateNoteService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// CreateNote create note info
func (s *CreateUserService) CreateNote(req *user.RegisterRequest) error {
	userModel := &db.User{
		UserName: req.UserName,
		Password: encryptPassWord(req.PassWord),
	}
	return db.CreateUser(s.ctx, []*db.User{userModel})
}

//TODO: 加密密码
func encryptPassWord(passWord string) string {
	return passWord
}
