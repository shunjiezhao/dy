package main

import (
	"first/pkg/constants"
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
		server.WithHostPorts(constants.ApiServerAddress),
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
		//relationGroup.GET("follow/list")
		relationGroup.GET("follower/list", follow.GetFollowerList(jwt.GetClaimsFromJWT))
		relationGroup.POST("action", follow.Follow(jwt.GetClaimsFromJWT))
	}

	r.Spin()
}
