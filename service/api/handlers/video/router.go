package video

import (
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/pkg/mq"
	"first/service/api/rpc/user"
	video2 "first/service/api/rpc/video"
	"github.com/gin-gonic/gin"
)

type (
	//Service 用户微服务代理
	Service struct {
		Video          video2.RpcProxyIFace
		User           user.RpcProxyIFace
		VideoPublisher []*mq.Publisher
	}
)

func NewVideo(face video2.RpcProxyIFace, userFace user.RpcProxyIFace,
	Pub []*mq.Publisher) *Service {
	if face == nil || userFace == nil {
		return nil
	}
	service := Service{
		Video:          face,
		User:           userFace,
		VideoPublisher: Pub,
	}

	return &service
}

func InitRouter(engine *gin.Engine) {

	_, jwtToken := middleware.JwtMiddle()

	// 创建publisher
	publishers := make([]*mq.Publisher, constants.VideoQCount)
	conn := mq.GetMqConnection()
	for i := 0; i < int(constants.VideoQCount); i++ {
		publishers[i] = mq.NewPublisher(conn, constants.SaveVideoExName,
			mq.GetSaveVideoQueueKey(int64(i)))
	}
	video := NewVideo(video2.NewVideoProxy(), user.NewUserProxy(), publishers)

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
