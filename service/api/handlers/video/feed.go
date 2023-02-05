package video

import (
	userPb "first/kitex_gen/user"
	videoPb "first/kitex_gen/video"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *Service) Feed(validate gin.HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param FeedRequest
			list  []*videoPb.Video
			users []*userPb.User
			req   *videoPb.GetVideoListRequest
		)

		err = c.ShouldBindQuery(&param)
		if param.LatestTime != 0 && err != nil {
			goto ParamErr

		}

		klog.Infof("[Feed]: 获取到参数 %#v", param)
		//	2. 获取数据 绑定

		if len(param.GetToken()) != 0 {
			validate(c)
			_, exists := c.Get(constants.IdentityKey)
			if !exists {
				handlers.SendResponse(c, errno.AuthorizationFailedErr) // token 验证失败
				goto ErrHandler

			}
			req = &videoPb.GetVideoListRequest{
				Uuid:      handlers.GetTokenUserId(c),
				TimeStamp: param.LatestTime,
			}

		} else {
			req = &videoPb.GetVideoListRequest{
				TimeStamp: param.LatestTime, // 获取当前点之后的
			}
		}
		// MS -> S
		if req.TimeStamp > time.Now().Unix() {
			req.TimeStamp /= 1000
		}

		list, users, err = s.GetVideosAndUsers(c, req)
		if err != nil {
			handlers.SendResponse(c, errno.ServiceErr)
			goto ErrHandler

		}
		klog.Infof("获取成功: %v", list)
		klog.Infof("获取成功users: %v", users)

		SendVideoListResponse(c, handlers.PackVideos(list, users, false), errno.Success)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	ErrHandler:
		c.Abort()
	}
}

//GetVideosAndUsers 获取视频 以及其用户信息
func (s *Service) GetVideosAndUsers(c *gin.Context, param *videoPb.GetVideoListRequest) ([]*videoPb.Video, []*userPb.User, error) {
	list, err := s.Video.GetVideoList(c, param)
	id := make([]int64, len(list))
	for i := 0; i < len(list); i++ {
		id[i] = list[i].Author
	}

	users, err := s.User.GetUsers(c, &userPb.GetUserSRequest{Id: id})
	if err != nil {

		klog.Errorf("[Video]: 获取视频用户信息失败: %v", err.Error())
		return nil, nil, err

	}
	return list, users, nil
}
