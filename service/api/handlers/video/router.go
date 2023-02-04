package video

import (
	"first/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {

	_, jwtToken := middleware.JwtMiddle()

	dy := engine.Group("/douyin")

	//	相关服务
	group := dy.Group("/publish")
	{
		//TODO:
		video := NewVideo(&defaultStorage{})
		group.Use(jwtToken)
		group.POST("/action/", video.Publish())
		group.GET("/list/", video.List())
	}

}
