package video

import (
	userPb "first/kitex_gen/user"
	videoPb "first/kitex_gen/video"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) List() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param ListRequest
			list  []*videoPb.Video

			user *userPb.User
		)

		err = c.ShouldBind(&param)
		if err != nil {
			klog.Errorf("[获取发布视频]: 绑定参数失败: %v", err.Error())
			goto ParamErr

		}

		klog.Infof("[获取发布视频]: 获取到 参数", param)
		//	2. 获取数据 绑定

		list, err = s.Video.GetVideoList(c, &videoPb.GetVideoListRequest{
			Author:     handlers.GetTokenUserId(c),
			GetAuthor_: true,
		})

		user, err = s.User.GetUserInfo(c, &userPb.GetUserRequest{Id: param.GetUserId()})
		if err != nil {
			klog.Errorf("[获取发布视频]: 获取视频用户信息失败: %v", err.Error())
			handlers.SendResponse(c, errno.ServiceErr)
			goto ErrHandler

		}

		SendVideoListResponse(c, handlers.PackVideos(list, []*userPb.User{user}, true), errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}
