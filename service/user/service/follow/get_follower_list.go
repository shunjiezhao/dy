package follow

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"first/service/user/pack"
	"log"
)

type GetFollowerUserListService struct {
	ctx context.Context
}

// NewGetFollowerUserListService new CreateNoteService
func NewGetFollowerUserListService(ctx context.Context) *GetFollowerUserListService {
	return &GetFollowerUserListService{ctx: ctx}
}

// GetFollowerUserList create note info
func (s *GetFollowerUserListService) GetFollowerUserList(req *user.GetFollowerListRequest) ([]*user.User, error) {
	log.Println("FollowerUserList: rpc 响应开始调用")
	list, err := db.GetFollowerUserList(s.ctx, req.Id)
	if err != nil {
		log.Println("获取粉丝列表失败: ", err.Error())
		return nil, err
	}

	return pack.Users(list), nil
}
