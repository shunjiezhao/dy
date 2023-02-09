package video

import (
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *Service) Feed(validate gin.HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param common.FeedRequest
			list  []*handlers.Video
		)

		err = c.ShouldBindQuery(&param)
		if param.LatestTime != 0 && err != nil {
			goto ParamErr

		}
		if len(param.GetToken()) != 0 {
			validate(c)
			uuid, exists := c.Get(constants.IdentityKey)
			if !exists {
				handlers.SendResponse(c, errno.AuthorizationFailedErr) // token 验证失败
				goto ErrHandler

			}
			param.Uuid = uuid.(int64)

		}

		klog.Infof("[Feed]: 获取到参数 %#v", param)

		// MS -> S
		if param.LatestTime > time.Now().Unix() {
			param.LatestTime /= 1000
		}
		if param.LatestTime == 0 {
			param.LatestTime = time.Now().Unix()
		}

		list, err = s.GetVideosAndUsers(c, &param, false)
		if err != nil {
			klog.Errorf("rpc 获取视频列表失败 %v", err)
			handlers.SendResponse(c, errno.RemoteErr)
			goto ErrHandler

		}
		klog.Infof("获取成功: %v", list)

		common.SendVideoListResponse(c, list, errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}

//GetVideosAndUsers 获取视频 以及其用户信息
func (s *Service) GetVideosAndUsers(c *gin.Context, param *common.FeedRequest, isOne bool) ([]*handlers.Video, error) {
	return s.Video.GetVideoList(c, param)
}
