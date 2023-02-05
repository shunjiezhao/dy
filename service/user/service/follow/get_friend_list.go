package follow

import (
	"context"
	"first/kitex_gen/user"
)

type GetFriendListService struct {
	ctx context.Context
}

// NewGetFriendListService new CreateNoteService
func NewGetFriendListService(ctx context.Context) *GetFollowUserListService {
	return &GetFollowUserListService{ctx: ctx}
}

// GetFriendLList create note info
func (s *GetFollowUserListService) GetFriendLList(req *user.GetFollowListRequest) ([]*user.User, error) {
	// TODO: 获取消息
	return nil, nil
}
