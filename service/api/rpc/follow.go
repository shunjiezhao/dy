package rpc

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
)

func GetFollowList(ctx context.Context, req *userPb.GetFollowListRequest) ([]*userPb.User, error) {
	resp, err := userClient.GetFollowList(ctx, req)
	if err != nil || resp.User == nil {
		return nil, err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
func GetFollowerList(ctx context.Context, req *userPb.GetFollowerListRequest) ([]*userPb.User, error) {
	resp, err := userClient.GetFollowerList(ctx, req)
	if err != nil || resp.User == nil {
		return nil, err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
func FollowUser(ctx context.Context, req *userPb.FollowRequest) error {
	resp, err := userClient.Follow(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}
func UnFollowUser(ctx context.Context, req *userPb.FollowRequest) error {
	resp, err := userClient.UnFollow(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}
