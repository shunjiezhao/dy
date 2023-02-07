package follow

import (
	"context"
	"first/kitex_gen/user"
	userDB "first/service/user/model/db"
	"first/service/user/pack"
)

type GetFriendListService struct {
	ctx context.Context
}

// NewGetFriendListService new CreateNoteService
func NewGetFriendListService(ctx context.Context) *GetFollowUserListService {
	return &GetFollowUserListService{ctx: ctx}
}

// GetFriendLList create note info
func (s *GetFollowUserListService) GetFriendLList(req *user.GetFriendRequest) ([]*user.FriendUser, error) {
	list, err := userDB.GetFriendChatList(userDB.DB, s.ctx, req.FromUserId)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(list))
	idx := make(map[int64]int, len(list))
	for i := 0; i < len(list); i++ {
		ids[i] = list[i].OtherId
		idx[ids[i]] = i // id -> list[]
	}

	users, err := userDB.MGetUsers(userDB.DB, s.ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make([]*user.FriendUser, len(list))
	for i := 0; i < len(users); i++ {
		res[i] = &user.FriendUser{}
		res[i].User = pack.User(users[i])
		j := idx[users[i].Uuid]
		if list[j].MySend {
			res[i].MsgType = 1
		}
		res[i].Message = list[j].Content
	}

	return res, nil
}
