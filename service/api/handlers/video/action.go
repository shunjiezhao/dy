package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
)

const defaultMaxSize int64 = 32 << 20

func (s *Service) Publish() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err        error
			act        PublishRequest
			fileHeader *multipart.FileHeader
			fileInfo   *FileInfo
		)
		// 1. 检查文件大小
		err = c.Request.ParseMultipartForm(defaultMaxSize)
		if err != nil {
			goto ParamErr

		}

		err = c.ShouldBind(&act)
		if err != nil {
			goto ParamErr

		}

		log.Println("获取到 参数", act)
		//	2. 获取数据 绑定

		err = c.ShouldBind(&act)
		if err != nil {
			log.Println("bind body error")
			return

		}

		fileHeader, err = c.FormFile("data") // 返回第一个
		if err != nil || fileHeader == nil {
			log.Println("获取文件头部error", err)
			goto ParamErr

		}
		//	3. 调用储存接口
		fileInfo, err = s.Storage.UploadFile(fileHeader)
		if err != nil {
			log.Printf("svc.UploadFile err: %v\n", err)
			handlers.BuildResponse(errno.NewErrNo(errno.ServiceErrCode, err.Error()))
			goto errHandler

		}

		log.Println("save file", fileInfo)
		handlers.BuildResponse(errno.Success)
		return

	ParamErr:
		handlers.BuildResponse(errno.ParamErr)
	errHandler:
		c.Abort()
	}
}
