package video

import (
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/service/api/handlers/storage"
	video2 "first/service/api/rpc/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {

	_, jwtToken := middleware.JwtMiddle()

	dy := engine.Group("/douyin")

	//	相关服务
	group := dy.Group("/publish")
	{
		factory := storage.DefaultOssFactory{
			Key: constants.OssSecretKey,
			Id:  constants.OssSecretID,
			Url: constants.OssUrl,
		}
		video := NewVideo(factory, video2.NewVideoProxy())
		group.Use(jwtToken)
		group.POST("/action/", video.Publish())
		group.GET("/list/", video.List())
	}

}
