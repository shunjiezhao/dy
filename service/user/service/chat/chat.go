package chat

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/util"
	userDB "first/service/user/model/db"
	"first/service/user/pack"
	"time"
)

type Service struct {
	ctx context.Context
}

// NewChatService new CreateNoteService
func NewChatService(ctx context.Context) *Service {
	return &Service{ctx: ctx}
}

// SaveChat 保存消息
func (s *Service) SaveChat(req *userPb.SaveMsgRequest) error {
	return userDB.SaveChat(userDB.DB, s.ctx, &userDB.Message{
		Id:           util.NextVal(),
		FromUserUuid: req.FromUserId,
		ToUserUuid:   req.ToUserId,
		Content:      req.Content,
		Base:         userDB.Base{CreatedAt: time.Unix(req.CreatedAtS, 0)},
	})
}

// GetChatList 获取消息列表
func (s *Service) GetChatList(req *userPb.GetChatListRequest) ([]*userPb.Message, error) {
	list, err := userDB.GetChatList(userDB.DB, s.ctx, req.FromUserId, req.ToUserId)
	return pack.Messages(list), err
}

// GetFriendChatList 获取好友的消息列表
func (s *Service) GetFriendChatList(req *userPb.GetFriendChatRequest) ([]*userPb.Message, error) {
	list, err := userDB.GetFriendChatList(userDB.DB, s.ctx, req.FromUserId)
	msgS := make([]*userPb.Message, len(list))
	for i := 0; i < len(list); i++ {
		msgS[i] = &userPb.Message{}
		if list[i].MySend {
			msgS[i].FromUserId = req.FromUserId
			msgS[i].ToUserId = list[i].UUid
		} else {
			msgS[i].ToUserId = req.FromUserId
			msgS[i].FromUserId = list[i].UUid
		}
		msgS[i].Content = list[i].Content
	}
	return msgS, err
}
