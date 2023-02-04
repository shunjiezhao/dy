package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/storage"
	"first/service/api/rpc/video"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

type (

	//Service 用户微服务代理
	Service struct {
		storage.Storage
		video.RpcProxyIFace
	}
)

func NewVideo(factory storage.StorageFactory, face video.RpcProxyIFace) *Service {
	service := Service{}
	if factory != nil {
		service.Storage = factory.Factory()
	}
	if face != nil {
		service.RpcProxyIFace = face
	}
	if factory == nil || face == nil {
		return nil
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
		Data []byte `json:"data"`
	}
	ListResponse struct {
		handlers.Response
	}
)

func SendPublishResponse(c *gin.Context, err error) {
	if err == nil {
		err = errno.Success
	}
	c.JSON(consts.StatusOK, PublishResponse{
		Response: handlers.BuildResponse(err),
	})
}
