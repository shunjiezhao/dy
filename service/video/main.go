package main

import (
	userService "first/kitex_gen/user/userservice"
	video "first/kitex_gen/video/videoservice"
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/pkg/util"
	"first/service/video/model"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func Init() {
	model.InitVideoDB()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	resolver, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	trace, _ := util.Trace("video-" + constants.UserServiceName)

	userClient, err := userService.NewClient(
		constants.UserServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithSuite(trace),
		client.WithResolver(resolver), // etcd
	)
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		panic(err)
	}
	Init()
	service := &VideoServiceImpl{}
	suite, close := util.SrvTrace(constants.VideoServiceName + "-service")
	defer close.Close()
	service.userRpc = userClient
	svr := video.NewServer(service,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.VideoServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                              // middleWare
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr), // address
		server.WithSuite(suite),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 10000}), // limit
		server.WithRegistry(r))

	clean := service.ConsumerStart()
	defer clean()
	clean = service.UpdateVideoInfoConStart()
	defer clean()
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
