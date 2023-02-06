package main

import (
	"context"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/pack"
	"first/service/user/service/chat"
	comment "first/service/user/service/comment"
	"first/service/user/service/follow"
	user2 "first/service/user/service/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
)

//TODO:
// 1.对于 string 类型 进行 SQL 注入检查
// 2.参数检查

type ChatServiceImpl struct{}

func (c *ChatServiceImpl) SendMsg(ctx context.Context, req *user.SaveMsgRequest) (resp *user.SaveMsgResponse, err error) {
	resp = new(user.SaveMsgResponse)
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

// isAdd 是否是发布评论
func isAdd(i int32) bool {
	if i == 1 {
		return true
	}
	return false
}
func (s *UserServiceImpl) ActionComment(ctx context.Context, req *user.ActionCommentRequest) (resp *user.
	ActionCommentResponse, err error) {
	resp = new(user.ActionCommentResponse)
	if isAdd(req.ActionType) { // 创建
		err = comment.NewCommentService(ctx).CreateComment(req)
	} else {
		err = comment.NewCommentService(ctx).DeleteComment(req)
	}

	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

func (s *UserServiceImpl) GetComment(ctx context.Context, req *user.GetCommentRequest) (resp *user.GetCommentResponse, err error) {
	resp = new(user.GetCommentResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		klog.Infof("[GetComment]: 参数有误")
		return resp, nil

	}
	resp.Comment, err = comment.NewCommentService(ctx).GetComment(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, req *user.GetUserSRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return
	}
	if req.Uuid == 0 {
		resp.User, err = user2.NewGetUserService(ctx).GetUserS(req)
	} else {
		resp.User, err = user2.NewGetUserService(ctx).GetUserSWithLogin(req) // 登陆用户获取朋友列表, 关注字段
	}
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserResponse, err error) {
	log.Println("user rpc server: get user")
	resp = new(user.GetUserResponse)
	resp.User, err = user2.NewGetUserService(ctx).GetUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	if len(req.UserName) == 0 || len(req.PassWord) == 0 {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return
	}
	resp = new(user.RegisterResponse)
	resp.Id, err = user2.NewCreateUserService(ctx).CreateUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse,
	err error) {
	log.Println("user rpc server: check user")
	resp = &user.CheckUserResponse{} // 使用 new 里面的resp 不会初始化
	resp.Id, err = user2.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.AuthorizationFailedErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.GetFollowerListRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	resp.User, err = follow.NewGetFollowerUserListService(ctx).GetFollowerUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.GetFollowListRequest) (resp *user.UserListResponse, err error) {
	resp = new(user.UserListResponse)
	resp.User, err = follow.NewGetFollowUserListService(ctx).GetFollowUserList(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// Follow implements the UserServiceImpl interface.
func (s *UserServiceImpl) Follow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse, err error) {
	log.Println("user rpc server: follow user")
	//??? 如果再次关注会怎么样?
	resp = new(user.FollowResponse)
	_, err = user2.NewFollowUserService(ctx).FollowUser(req)
	resp.Resp = pack.BuildBaseResp(err)

	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// UnFollow implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnFollow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse,
	err error) {
	log.Println("user rpc server: follow user")
	resp = new(user.FollowResponse)
	_, err = user2.NewUnFollowUserService(ctx).UnFollowUser(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(err)
		return nil, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}

// GetFriendList 返回好友列表
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.FollowRequest) (resp *user.GetFriendResponse, err error) {

	return
}
