package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"time"
)

const defaultMaxSize int64 = 32 << 20

func (s *Service) Publish() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err        error
			param      PublishRequest
			fileHeader *multipart.FileHeader
			//fileInfo   storage.AccessUrl
		)
		// 1. 检查文件大小
		err = c.Request.ParseMultipartForm(defaultMaxSize)
		if err != nil {
			goto ParamErr

		}

		err = c.ShouldBind(&param)
		if err != nil {
			goto ParamErr

		}

		klog.Infof("[发布视频]: 获取到 参数:%v", param)
		//	2. 获取数据 绑定

		fileHeader, err = c.FormFile("data") // 返回第一个
		if err != nil || fileHeader == nil {
			klog.Errorf("[发布视频]: 获取文件头部error", err)
			goto ParamErr

		}
		//	3. 调用储存接口
		s.Storage.UploadFile(param.Title, fileHeader, handlers.GetTokenUserId(c), time.Now())
		if err != nil {
			klog.Errorf("[发布视频]:  rpc 出现错误: %v\n", err)
			handlers.SendResponse(c, errno.NewErrNo(errno.ServiceErrCode, err.Error()))
			goto errHandler

		}
		klog.Info("[发布视频]: 	成功")

		handlers.SendResponse(c, errno.Success)
		return

	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	errHandler:
		c.Abort()
	}
}
