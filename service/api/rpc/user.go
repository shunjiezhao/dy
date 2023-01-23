package rpc

import (
	"context"
	userPb "first/kitex_gen/user"
	userService "first/kitex_gen/user/userservice"
	"first/pkg/constants"
	"first/pkg/middleware"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

// userservice rpc client
var userClient userService.Client

func initUserRpc() {
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

// Hello 是对rpc调用的包装
func Register(ctx context.Context, req *userPb.RegisterRequest) (int64, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return -1, err
	}
	return resp.Id, nil
}
