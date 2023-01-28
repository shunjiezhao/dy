package main

import (
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/service/api/handlers/follow"
	"first/service/api/handlers/user"
	"first/service/api/rpc"
	"github.com/gin-gonic/gin"
)

func Init() {
	rpc.InitRPC()
}
func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	Init()
	jwt, jwtToken := middleware.JwtMiddle()
	engine := gin.Default()

	dy := engine.Group("/douyin")

	// 用户相关
	{
		dy.GET("/user/", jwtToken, user.GetInfo())
		dy.POST("/user/register/", user.Register(jwt.TokenGenerator, nil))
		dy.POST("/user/login/", user.Login(), jwt.LoginHandler)
	}
	// 社交接口的相关实现
	relationGroup := dy.Group("relation")
	{
		relationGroup.Use(jwtToken)
		followSrv := follow.New(nil)
		relationGroup.GET("/follow/list/", followSrv.GetFollowList())
		relationGroup.GET("/follower/list/", followSrv.GetFollowerList())
		relationGroup.POST("/action/", followSrv.Follow())
	}

	engine.Run(constants.ApiServerAddress)
}
