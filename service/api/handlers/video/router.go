package video

import (
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/service/api/handlers/storage"
	"first/service/api/rpc/user"
	video2 "first/service/api/rpc/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {

	_, jwtToken := middleware.JwtMiddle()
	factory := storage.DefaultOssFactory{
		Key: constants.OssSecretKey,
		Id:  constants.OssSecretID,
		Url: constants.OssUrl,
	}
	video := NewVideo(factory, video2.NewVideoProxy(), user.NewUserProxy())

	dy := engine.Group("/douyin")
	dy.GET("/feed/", video.Feed(jwtToken)) // 获取视频流

	//	相关服务
	group := dy.Group("/publish")
	{

		group.Use(jwtToken)
		group.POST("/action/", video.Publish()) // 发布视频
		group.GET("/list/", video.List())       // 获取发布的视频
	}

	favourite := dy.Group("/favorite")
	{
		favourite.Use(jwtToken)
		favourite.POST("/action/", video.Like())   // 喜欢
		favourite.GET("/list/", video.LikeVideo()) // 喜欢的视频列表
	}
}
