package video

import (
	"bytes"
	"encoding/gob"
	"first/pkg/errno"
	"first/pkg/mq"
	"first/pkg/storage"
	"first/pkg/util"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"time"
)

const defaultMaxSize int64 = 32 << 20

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func (s *Service) Publish() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err        error
			param      common.PublishRequest
			fileHeader *multipart.FileHeader
			//fileInfo   storage.AccessUrl
			file     multipart.File
			data     []byte
			uuid     int64
			saveInfo *storage.Info
			ctx      = c.Request.Context()
			byte2    bytes.Buffer
			encoder  *gob.Encoder
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

		//3.压缩文件
		file, err = fileHeader.Open()
		defer file.Close()
		if err != nil {
			return
		}

		data, err = io.ReadAll(file)
		if err != nil {
			return
		}

		saveInfo = &storage.Info{
			Data:  data,
			Time:  time.Now().Unix(),
			Uuid:  handlers.GetTokenUserId(c),
			Title: param.Title,
		}
		encoder = gob.NewEncoder(&byte2)
		err = encoder.Encode(saveInfo)

		if err != nil {
			klog.Error("转换失败")
			handlers.SendResponse(c, err)
			goto errHandler
		}
		data, err = util.Compress(byte2.Bytes())
		if err != nil {
			klog.Error("加密失败")
			handlers.SendResponse(c, err)
			goto errHandler

		}
		//4.传入消息队列
		err = s.Publisher[mq.GetSaveVideoIdx(uuid)].Publish(ctx, data)
		if err != nil {
			klog.Errorf("[发布视频]:  发送消息队列失败 %v", err)
			goto errHandler

		}

		//klog.Info("[发布视频]: 	成功")
		handlers.SendResponse(c, errno.Success)
		return

	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	errHandler:
		c.Abort()
	}
}
