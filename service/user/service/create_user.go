package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
)

var ss *snowflake.Snowflake

func init() {
	var err error
	ss, err = snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		panic(err)
	}
}

type CreateUserService struct {
	ctx context.Context
}

// NewCreateNoteService new CreateNoteService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// CreateUser create note info
func (s *CreateUserService) CreateUser(req *user.RegisterRequest) (int64, error) {
	println("rpc 响应开始调用")
	userModel := &db.User{
		Uuid:     ss.NextVal(),
		UserName: req.UserName,
		Password: encryptPassWord(req.PassWord),
	}
	return db.CreateUser(s.ctx, []*db.User{userModel})
}

//TODO: 加密密码
func encryptPassWord(passWord string) string {
	return passWord
}
