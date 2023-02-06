package user

import (
	"context"
	userPb "first/kitex_gen/user"
	userService "first/kitex_gen/user/userservice"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

// userClient rpc client
var userClient userService.Client

//go:generate mockgen -destination=../mock/male_mock.go -package=mock first/service/api/rpc/user RpcProxyIFace
type RpcProxyIFace interface {
	Register(ctx context.Context, req *userPb.RegisterRequest) (int64, error)
	GetUserInfo(ctx context.Context, req *userPb.GetUserRequest) (*userPb.User, error)
	CheckUser(ctx context.Context, req *userPb.CheckUserRequest) (int64, error)
	UnFollowUser(ctx context.Context, req *userPb.FollowRequest) error
	FollowUser(ctx context.Context, req *userPb.FollowRequest) error
	GetFollowerList(ctx context.Context, req *userPb.GetFollowerListRequest) ([]*userPb.User, error)
	GetFollowList(ctx context.Context, req *userPb.GetFollowListRequest) ([]*userPb.User, error)
	GetUsers(ctx context.Context, Req *userPb.GetUserSRequest) ([]*userPb.User, error)

	ActionComment(ctx context.Context, Req *userPb.ActionCommentRequest) (r *userPb.Comment, err error) //评论操作
	GetComment(ctx context.Context, Req *userPb.GetCommentRequest) (r []*userPb.Comment, err error)     // 获取评论
}

type RpcProxy struct {
	userClient userService.Client
}

func NewUserProxy() RpcProxyIFace {
	return &RpcProxy{userClient: userClient}
}

func InitApiUserRpc() {
	var err error
	resolver, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	userClient, err = userService.NewClient(
		constants.UserServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithResolver(resolver), // etcd
	)
	if err != nil {
		panic(err)
	}
}
func respIsErr(Resp *userPb.BaseResp) bool {
	return Resp != nil && Resp.StatusCode != errno.SuccessCode
}

// Register rpc调用, 如果成功返回 userid
func (proxy RpcProxy) Register(ctx context.Context, req *userPb.RegisterRequest) (int64, error) {
	resp, err := proxy.userClient.Register(ctx, req)
	if err != nil {
		return 0, errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return 0, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.Id, nil
}

//CheckUser rpc调用, 检查用户是否存在,如果存在返回 userid
func (proxy RpcProxy) CheckUser(ctx context.Context, req *userPb.CheckUserRequest) (int64, error) {
	resp, err := proxy.userClient.CheckUser(ctx, req)
	if err != nil {
		return 0, errno.RemoteErr
	}
	// NOTICE: 注意判断, 可能上方用 new 导致 null pointer 异常
	if respIsErr(resp.Resp) {
		return 0, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.Id, nil
}
func (proxy RpcProxy) GetUserInfo(ctx context.Context, req *userPb.GetUserRequest) (*userPb.User, error) {
	resp, err := proxy.userClient.GetUser(ctx, req)
	if err != nil || resp.User == nil {
		return nil, errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
func (proxy RpcProxy) GetUsers(ctx context.Context, Req *userPb.GetUserSRequest) ([]*userPb.User, error) {
	if len(Req.Id) == 0 {
		return nil, errno.RemoteErr
	}

	resp, err := proxy.userClient.GetUsers(ctx, Req)
	if err != nil || resp.User == nil {
		return nil, errno.RemoteErr
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return resp.User, nil
}
