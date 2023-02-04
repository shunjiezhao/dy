package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type FileInfo struct {
	AccessUrl string
}

type Storage interface {
	UploadFile(*multipart.FileHeader) (*FileInfo, error) // 返回我们的 获取链接
}

//Service 用户微服务代理
type Service struct {
	Storage
}

func NewVideo(storage Storage) *Service {
	return &Service{storage}
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

func SendRegisterResponse(c *gin.Context) {
	c.JSON(consts.StatusOK, PublishResponse{
		Response: handlers.BuildResponse(errno.Success),
	})
}
