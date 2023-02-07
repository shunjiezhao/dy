package main

import (
	video "first/kitex_gen/video/videoservice"
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/service/video/model"
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
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		panic(err)
	}
	Init()
	service := new(VideoServiceImpl)
	svr := video.NewServer(service,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.VideoServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                              // middleWare
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),                                         // address
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}), // limit
		server.WithRegistry(r))

	clean := service.ConsumerStart()
	defer clean()
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
