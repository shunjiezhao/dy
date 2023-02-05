package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
)

func (proxy RpcProxy) GetFollowList(ctx context.Context, req *userPb.GetFollowListRequest) ([]*userPb.User, error) {
	resp, err := proxy.userClient.GetFollowList(ctx, req)
	if err != nil || resp.User == nil {
		return nil, err
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
func (proxy RpcProxy) GetFollowerList(ctx context.Context, req *userPb.GetFollowerListRequest) ([]*userPb.User, error) {
	resp, err := proxy.userClient.GetFollowerList(ctx, req)
	if err != nil || resp.User == nil {
		return nil, err
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
func (proxy RpcProxy) FollowUser(ctx context.Context, req *userPb.FollowRequest) error {
	resp, err := proxy.userClient.Follow(ctx, req)
	if err != nil {
		return err
	}
	if respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}
func (proxy RpcProxy) UnFollowUser(ctx context.Context, req *userPb.FollowRequest) error {
	resp, err := proxy.userClient.UnFollow(ctx, req)
	if err != nil {
		return err
	}
	if respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}
