package rpc

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

// userservice rpc client
var userClient userService.Client

func initApiUserRpc() {
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

// Register rpc调用, 如果成功返回 userid
func Register(ctx context.Context, req *userPb.RegisterRequest) (int64, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return -1, err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return 0, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.Id, nil
}

//CheckUser rpc调用, 检查用户是否存在,如果存在返回 userid
func CheckUser(ctx context.Context, req *userPb.CheckUserRequest) (int64, error) {
	resp, err := userClient.CheckUser(ctx, req)
	if err != nil {
		return -1, err
	}
	// NOTICE: 注意判断, 可能上方用 new 导致 null pointer 异常
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return 0, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.Id, nil
}
func GetUserInfo(ctx context.Context, req *userPb.GetUserRequest) (*userPb.User, error) {
	resp, err := userClient.GetUser(ctx, req)
	if err != nil || resp.User == nil {
		return nil, err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.User, nil
}
