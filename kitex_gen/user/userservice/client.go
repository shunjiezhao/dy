// Code generated by Kitex v0.4.4. DO NOT EDIT.

package userservice

import (
	"context"
	user "first/kitex_gen/user"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Register(ctx context.Context, Req *user.RegisterRequest, callOptions ...callopt.Option) (r *user.RegisterResponse, err error)
	CheckUser(ctx context.Context, Req *user.CheckUserRequest, callOptions ...callopt.Option) (r *user.CheckUserResponse, err error)
	GetUser(ctx context.Context, Req *user.GetUserRequest, callOptions ...callopt.Option) (r *user.GetUserResponse, err error)
	GetFollowerList(ctx context.Context, Req *user.GetFollowerListRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error)
	GetFollowList(ctx context.Context, Req *user.GetFollowListRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error)
	Follow(ctx context.Context, Req *user.FollowRequest, callOptions ...callopt.Option) (r *user.FollowResponse, err error)
	UnFollow(ctx context.Context, Req *user.FollowRequest, callOptions ...callopt.Option) (r *user.FollowResponse, err error)
	GetFriendList(ctx context.Context, Req *user.GetFriendRequest, callOptions ...callopt.Option) (r *user.GetFriendResponse, err error)
	GetUsers(ctx context.Context, Req *user.GetUserSRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error)
	ActionComment(ctx context.Context, Req *user.ActionCommentRequest, callOptions ...callopt.Option) (r *user.ActionCommentResponse, err error)
	GetComment(ctx context.Context, Req *user.GetCommentRequest, callOptions ...callopt.Option) (r *user.GetCommentResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kUserServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kUserServiceClient struct {
	*kClient
}

func (p *kUserServiceClient) Register(ctx context.Context, Req *user.RegisterRequest, callOptions ...callopt.Option) (r *user.RegisterResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Register(ctx, Req)
}

func (p *kUserServiceClient) CheckUser(ctx context.Context, Req *user.CheckUserRequest, callOptions ...callopt.Option) (r *user.CheckUserResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CheckUser(ctx, Req)
}

func (p *kUserServiceClient) GetUser(ctx context.Context, Req *user.GetUserRequest, callOptions ...callopt.Option) (r *user.GetUserResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUser(ctx, Req)
}

func (p *kUserServiceClient) GetFollowerList(ctx context.Context, Req *user.GetFollowerListRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowerList(ctx, Req)
}

func (p *kUserServiceClient) GetFollowList(ctx context.Context, Req *user.GetFollowListRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFollowList(ctx, Req)
}

func (p *kUserServiceClient) Follow(ctx context.Context, Req *user.FollowRequest, callOptions ...callopt.Option) (r *user.FollowResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Follow(ctx, Req)
}

func (p *kUserServiceClient) UnFollow(ctx context.Context, Req *user.FollowRequest, callOptions ...callopt.Option) (r *user.FollowResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UnFollow(ctx, Req)
}

func (p *kUserServiceClient) GetFriendList(ctx context.Context, Req *user.GetFriendRequest, callOptions ...callopt.Option) (r *user.GetFriendResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFriendList(ctx, Req)
}

func (p *kUserServiceClient) GetUsers(ctx context.Context, Req *user.GetUserSRequest, callOptions ...callopt.Option) (r *user.UserListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUsers(ctx, Req)
}

func (p *kUserServiceClient) ActionComment(ctx context.Context, Req *user.ActionCommentRequest, callOptions ...callopt.Option) (r *user.ActionCommentResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ActionComment(ctx, Req)
}

func (p *kUserServiceClient) GetComment(ctx context.Context, Req *user.GetCommentRequest, callOptions ...callopt.Option) (r *user.GetCommentResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetComment(ctx, Req)
}
