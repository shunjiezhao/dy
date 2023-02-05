package user

import (
	"first/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine, UserService *Service) {

	jwt, jwtToken := middleware.JwtMiddle()

	dy := engine.Group("/douyin")
	user := dy.Group("/user/")

	// 用户相关
	{
		user.GET("", jwtToken, UserService.GetInfo())
		user.POST("/register/", UserService.Register(jwt.TokenGenerator))
		user.POST("/login/", UserService.Login(), jwt.LoginHandler)
	}
	// 社交接口的相关实现
	relationGroup := dy.Group("relation")
	{
		relationGroup.Use(jwtToken)
		relationGroup.GET("/follow/list/", UserService.GetFollowList())
		// 粉丝列表
		relationGroup.GET("/follower/list/", UserService.GetFollowerList())
		// 好友列表
		relationGroup.GET("/friend/list/", UserService.GetFriendList())
		relationGroup.POST("/action/", UserService.Follow())
	}

	// 评论
	comment := dy.Group("/comment")
	{
		comment.Use(jwtToken)
		comment.GET("/list/", UserService.GetCommentList())
		comment.GET("/action/", UserService.ActionComment())
	}
}
