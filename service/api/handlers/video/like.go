package video

import (
	videoPb "first/kitex_gen/video"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) LikeVideo() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param common.ListRequest
			list  []*handlers.Video
		)

		err = c.ShouldBind(&param)
		if err != nil {
			err = c.ShouldBindQuery(&param)
			if err != nil {
				goto ParamErr

			}
		}
		klog.Infof("[LikeVideo]: 获取到参数 %#v", param)

		list, err = s.GetVideosAndUsers(c, &common.FeedRequest{
			Uuid:   param.GetUserId(),
			IsLike: true,
		}, false)

		if err != nil {
			handlers.SendResponse(c, err)
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
func (s *Service) Like() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param common.FavoriteActionRequest
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
			handlers.SendResponse(c, err)
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
