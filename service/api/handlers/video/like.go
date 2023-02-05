package video

import (
	userPb "first/kitex_gen/user"
	videoPb "first/kitex_gen/video"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) LikeVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param ListRequest
			list  []*videoPb.Video
			users []*userPb.User
			req   *videoPb.GetVideoListRequest
		)

		err = c.ShouldBind(&param)
		if err != nil {
			err = c.ShouldBindQuery(&param)
			if err != nil {
				goto ParamErr

			}
		}
		klog.Infof("[LikeVideo]: 获取到参数 %#v", param)

		req = &videoPb.GetVideoListRequest{
			Author:     handlers.GetTokenUserId(c),
			GetAuthor_: true,
		}
		list, users, err = s.GetVideosAndUsers(c, req)
		if err != nil {
			handlers.SendResponse(c, errno.ServiceErr)
			goto ErrHandler

		}

		SendVideoListResponse(c, handlers.PackVideos(list, users, false), errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}
func (s *Service) Like() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param FavoriteActionRequest
			req   *videoPb.LikeVideoRequest
		)

		err = c.ShouldBind(&param)
		if err != nil {
			err = c.ShouldBindQuery(&param)
			if err != nil {
				goto ParamErr

			}
		}
		klog.Infof("[Like]: 获取到参数 %#v", param)

		req = &videoPb.LikeVideoRequest{
			VideoId:    param.VideoId,
			ActionType: &videoPb.LikeVideoAction{ActionType: int32(param.VideoFavActionType)},
			Uuid:       handlers.GetTokenUserId(c),
		}
		err = s.Video.LikeVideo(c, req)
		if err != nil {
			handlers.SendResponse(c, errno.ServiceErr)
			goto ErrHandler

		}

		handlers.SendResponse(c, errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}
