package user

import (
	"first/pkg/constants"
	"first/pkg/middleware"
	"first/pkg/mq"
	user2 "first/service/api/rpc/user"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

//Service 用户微服务代理
type Service struct {
	rpc       user2.RpcProxyIFace
	chatSrv   user2.ChatProxy
	publisher map[string][]*mq.Publisher
	redis     *redis.Client
}

func New(rpc user2.RpcProxyIFace, charSrv user2.ChatProxy) *Service {
	if rpc == nil {
		rpc = user2.NewUserProxy()
	}
	if charSrv == nil {
		charSrv = user2.NewChatRpcProxy()
	}

	return &Service{
		rpc:       rpc,
		chatSrv:   charSrv,
		publisher: createPublisher(),
	}
}
func createPublisher() map[string][]*mq.Publisher {
	m := map[string][]*mq.Publisher{}
	publishers := make([]*mq.Publisher, constants.UActionCommentQCount)
	conn := mq.GetMqConnection()
	for i := 0; i < int(constants.UActionCommentQCount); i++ {
		publishers[i] = mq.NewPublisher(conn, constants.UActionCommentExName, mq.UGetActionCommentQueueKey(int64(i)))
	}
	m[constants.UActionCommentKey] = publishers

	publishers = make([]*mq.Publisher, constants.VActionVideoComCountQCount)
	for i := 0; i < int(constants.VActionVideoComCountQCount); i++ {
		publishers[i] = mq.NewPublisher(conn, constants.VActionVideoComCountExName, mq.VGetActionVideoComCountQueueKey(int64(i)))
	}
	m[constants.VActionVideoComCountKey] = publishers

	return m
}

func InitRouter(engine *gin.Engine, UserService *Service) {

	jwt, jwtToken := middleware.JwtMiddle()

	dy := engine.Group("/douyin")
	user := dy.Group("/user/")

	// 用户相关
	{
		user.GET("", jwtToken, UserService.GetInfo())
		user.POST("/register/", UserService.Register(), jwt.LoginHandler)
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
		comment.POST("/action/", UserService.ActionComment())
	}
	message := dy.Group("/message/")
	{
		message.Use(jwtToken)
		message.GET("/chat/", UserService.GetChatList())
		message.POST("/action/", UserService.SendMsg())
	}
}
