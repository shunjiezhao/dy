package handler

import (
	"context"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/pack"
	"first/service/user/service/chat"
	"first/service/user/service/follow"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ChatServiceImpl struct{}

func (c *ChatServiceImpl) GetFriendChatList(ctx context.Context, req *user.GetFriendChatRequest) (resp *user.GetFriendChatResponse, err error) {
	resp = new(user.GetFriendChatResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	resp.Msg, err = chat.NewChatService(ctx).GetFriendChatList(req)
	if err != nil {
		klog.Errorf("[获取好友消息记录] :获取出错 %v", err)
		resp.Resp = pack.BuildBaseResp(errno.ServiceErr)
		return
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)

	return
}

func (c *ChatServiceImpl) SendMsg(ctx context.Context, req *user.SaveMsgRequest) (resp *user.SaveMsgResponse, err error) {
	resp = new(user.SaveMsgResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	err = chat.NewChatService(ctx).SaveChat(req)
	if err != nil {
		klog.Errorf("[User.SendMsg]: msg保存失败 %v", err)
		resp.Resp = pack.BuildBaseResp(errno.MsgSaveErr)
		return
	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

func (c *ChatServiceImpl) GetChatList(ctx context.Context, req *user.GetChatListRequest) (resp *user.GetChatListResponse, err error) {
	resp = new(user.GetChatListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	resp.MessageList, err = chat.NewChatService(ctx).GetChatList(req)
	if err != nil {
		klog.Errorf("[User.GetChatList]: 获取消息失败 %v", err)
		resp.Resp = pack.BuildBaseResp(errno.MsgSaveErr)
		return
	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetFriendList 返回好友列表
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.GetFriendRequest) (resp *user.GetFriendResponse, err error) {
	resp = new(user.GetFriendResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	resp.UserList, err = follow.NewGetFriendListService(ctx).GetFriendLList(req)
	if err != nil {
		klog.Errorf("[获取好友消息记录] :获取出错 %v", err)
		resp.Resp = pack.BuildBaseResp(errno.ServiceErr)
		return
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
