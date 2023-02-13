package main

import (
	"first/kitex_gen/user/chatservice"
	user "first/kitex_gen/user/userservice"
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/pkg/util"
	"first/service/user/handler"
	"first/service/user/model"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	"sync"
)

func Init() {
	model.InitUserDB()
}
func Run(r registry.Registry) {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		panic(err)
	}
	impl := &handler.UserServiceImpl{}

	svr := user.NewServer(impl,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.UserServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                             // middleWare
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),                                         // address
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}), // limit
		server.WithRegistry(r))

	impl.UpdateVideoInfoConStart()
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	Init()
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		defer group.Done()
		Run(r) // user service
	}()

	// chat service
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		panic(err)
	}
	trace, _ := util.SrvTrace(constants.ChatServiceName + "-service")

	svr := chatservice.NewServer(&handler.ChatServiceImpl{},
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.ChatServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                             // middleWare
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr), // address
		server.WithSuite(trace),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}), // limit
		server.WithRegistry(r))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

	group.Wait()
}
