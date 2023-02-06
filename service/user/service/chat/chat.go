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
