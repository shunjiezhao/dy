package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/storage"
	"first/service/api/rpc/user"
	"first/service/api/rpc/video"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

type (

	//Service 用户微服务代理
	Service struct {
		storage.Storage
		Video video.RpcProxyIFace
		User  user.RpcProxyIFace
	}
)

func NewVideo(factory storage.StorageFactory, face video.RpcProxyIFace, userFace user.RpcProxyIFace) *Service {
	if factory == nil || face == nil || userFace == nil {
		return nil
	}
	service := Service{
		Storage: factory.Factory(),
		Video:   face,
		User:    userFace,
	}

	return &service
}

type (
	PublishRequest struct {
		handlers.Token
		Title string `json:"title" form:"title"`
	}
	PublishResponse struct {
		handlers.Response
	}
	ListRequest struct {
		handlers.Token
		handlers.UserId
	}
	ListResponse struct {
		handlers.Response
		VideoList []*handlers.Video `json:"video_list"`
	}

	FeedRequest struct {
		handlers.Token
		LatestTime int64 `json:"latest_time" form:"latest_time"` // 这个时间点以前的视频 [时间戳]
	}
	FeedResponse struct {
		handlers.Response
		VideoList []*handlers.Video `json:"video_list"`
	}
)
type (
	VideoFavActionType int32
	// Favourite

	FavoriteActionRequest struct {
		Token              string               `form:"token" `
		VideoId            int64                `form:"video_id"`
		VideoFavActionType `form:"action_type"` // 1-点赞，2-取消点赞
	}

	FavoriteActionResponse struct {
		StatusCode int32  `form:"status_code"`
		StatusMsg  string `form:"status_msg"`
	}
)

func (l VideoFavActionType) IsLike() bool {
	return l == 1
}
func SendPublishResponse(c *gin.Context, err error) {
	if err == nil {
		err = errno.Success
	}
	c.JSON(consts.StatusOK, PublishResponse{
		Response: handlers.BuildResponse(err),
	})
}

func SendVideoListResponse(c *gin.Context, list []*handlers.Video, err error) {
	c.JSON(consts.StatusOK, ListResponse{
		Response:  handlers.BuildResponse(err),
		VideoList: list,
	})
}
