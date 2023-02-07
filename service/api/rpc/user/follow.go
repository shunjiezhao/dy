package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	pack2 "first/service/api/rpc/user/pack"
	"github.com/cloudwego/kitex/pkg/klog"
)

func (proxy RpcProxy) GetFollowList(ctx context.Context, param *common.GetUserFollowListRequest) ([]*handlers.User, error) {

	req := &userPb.GetFollowListRequest{
		Id: param.GetUserId(),
	}
	resp, err := proxy.userClient.GetFollowList(ctx, req)
	if err != nil {
		klog.Errorf("[UserRpc.GetFollowList]: 失败 :%v", err)
		return nil, errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return pack2.Users(resp.User), nil
}
func (proxy RpcProxy) GetFollowerList(ctx context.Context, param *common.GetUserFollowerListRequest) ([]*handlers.User, error) {
	req := &userPb.GetFollowerListRequest{
		Id: param.GetUserId(),
	}

	resp, err := proxy.userClient.GetFollowerList(ctx, req)
	if err != nil {
		klog.Errorf("[UserRpc.GetFollowerList]: 失败 :%v", err)
		return nil, errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return pack2.Users(resp.User), nil
}

func (proxy RpcProxy) ActionFollow(ctx context.Context, param *common.ActionRequest) error {
	var (
		resp *userPb.FollowResponse
		err  error
	)

	req := &userPb.FollowRequest{
		FromUserId: param.FromUserId,
		ToUserId:   param.ToUserId,
	}

	if param.IsFollow() {
		resp, err = proxy.userClient.Follow(ctx, req)
	} else {
		resp, err = proxy.userClient.UnFollow(ctx, req)
	}

	if err != nil {
		klog.Errorf("[UserRpc.Action]: 失败 :%v", err)
		return errno.RemoteErr
	}

	if respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil

}
func (proxy RpcProxy) GetFriendList(ctx context.Context, param *common.FriendListRequest) ([]*handlers.FriendUser, error) {
	var (
		resp *userPb.GetFriendResponse
		err  error
	)

	req := &userPb.GetFriendRequest{
		FromUserId: param.GetUserId(),
	}
	resp, err = proxy.userClient.GetFriendList(ctx, req)
	if err != nil {
		klog.Errorf("[UserRpc.GetFriendList]: 失败 :%v", err)
		return nil, err
	}
	klog.Infof("[UserRpc.GetFriendList]: result: %v", resp.UserList)
	return pack2.FriendUsers(resp.UserList), nil
}
