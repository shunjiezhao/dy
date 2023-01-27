package main

import (
	"first/pkg/middleware"
	"first/service/api/handlers/follow"
	"first/service/api/handlers/user"
	"first/service/api/rpc"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Init() {
	rpc.InitRPC()
}
func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	Init()
	r := server.New(
		server.WithHostPorts(":8888"),
		server.WithHandleMethodNotAllowed(true),
	)
	jwt := middleware.JwtMiddle()

	dy := r.Group("/douyin")

	// 用户相关
	userGroup := dy.Group("user")
	{
		userGroup.GET("", user.GetInfo(jwt))
		userGroup.POST("register", user.Register(jwt.TokenGenerator, nil))
		userGroup.GET("login", user.Login(), jwt.LoginHandler)
	}
	// 社交接口的相关实现
	relationGroup := dy.Group("relation")
	{
		followSrv := follow.New(jwt.GetClaimsFromJWT)
		relationGroup.GET("follow/list", followSrv.GetFollowList())
		relationGroup.GET("follower/list", followSrv.GetFollowerList())
		relationGroup.POST("action", followSrv.Follow())
	}

	r.Spin()
}
