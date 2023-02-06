package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) List() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param common.ListRequest
			list  []*handlers.Video
		)

		err = c.ShouldBind(&param)
		if err != nil {
			klog.Errorf("[获取发布视频]: 绑定参数失败: %v", err.Error())
			goto ParamErr

		}

		klog.Infof("[获取发布视频]: 获取到 参数", param)
		//	2. 获取数据 绑定

		list, err = s.GetVideosAndUsers(c, &common.FeedRequest{
			Token:     param.Token,
			Author:    param.GetUserId(),
			GetAuthor: true,
		}, true)
		if err != nil {
			klog.Errorf("[获取发布视频]: 获取视频用户信息失败: %v", err.Error())
			handlers.SendResponse(c, errno.ServiceErr)
			goto ErrHandler

		}

		common.SendVideoListResponse(c, list, errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}
