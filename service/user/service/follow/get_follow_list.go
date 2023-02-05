package follow

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"first/service/user/pack"
	"log"
)

type GetFollowUserListService struct {
	ctx context.Context
}

// NewGetFollowUserListService new CreateNoteService
func NewGetFollowUserListService(ctx context.Context) *GetFollowUserListService {
	return &GetFollowUserListService{ctx: ctx}
}

// GetFollowUserList create note info
func (s *GetFollowUserListService) GetFollowUserList(req *user.GetFollowListRequest) ([]*user.User, error) {
	log.Println("FollowUserList: rpc 响应开始调用")
	list, err := db.GetFollowUserList(s.ctx, req.Id)
	if err != nil {
		log.Println("获取关注列表失败: ", err.Error())
		return nil, err
	}
	return pack.Users(list), nil
}
