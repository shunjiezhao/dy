package user

import (
	"first/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine, UserService *Service) {

	jwt, jwtToken := middleware.JwtMiddle()

	dy := engine.Group("/douyin")

	// 用户相关
	{
		dy.GET("/user/", jwtToken, UserService.GetInfo())
		dy.POST("/user/register/", UserService.Register(jwt.TokenGenerator))
		dy.POST("/user/login/", UserService.Login(), jwt.LoginHandler)
	}
	// 社交接口的相关实现
	relationGroup := dy.Group("relation")
	{
		relationGroup.Use(jwtToken)
		relationGroup.GET("/follow/list/", UserService.GetFollowList())
		// 好友列表就是粉丝列表
		relationGroup.GET("/friend/list/", UserService.GetFollowerList())
		relationGroup.GET("/follower/list/", UserService.GetFollowerList())
		relationGroup.POST("/action/", UserService.Follow())
	}
}
