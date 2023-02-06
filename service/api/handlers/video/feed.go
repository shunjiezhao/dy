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

		list, err = s.GetVideosAndUsers(c, &param, false)
		if err != nil {
			handlers.SendResponse(c, errno.ServiceErr)
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
	list, err := s.Video.GetVideoList(c, param)
	if err != nil {
		klog.Errorf("[Video]: 获取视频失败: %v", err.Error())
		return nil, err

	}
	if len(list) == 0 {
		return nil, nil
	}
	var id []int64
	if !isOne {
		id = make([]int64, len(list))
		for i := 0; i < len(list); i++ {
			id[i] = list[i].Author.Id
		}
	} else {
		id = make([]int64, 1)
		id[0] = list[0].Author.Id
	}

	users, err := s.User.GetUsers(c, &common.GetUserSRequest{Id: id, CurUserId: param.Uuid})
	if err != nil {
		klog.Errorf("[Video]: 获取视频用户信息失败: %v", err.Error())
		return nil, err

	}
	klog.Infof("获取成功users: %v", users)

	return Videos(list, users, isOne), nil
}
func Videos(videos []*handlers.Video, users []*handlers.User, isOne bool) []*handlers.Video {
	var (
		one *handlers.User
		idx map[int64]int
	)
	if isOne {
		one = users[0]
	} else {
		idx = make(map[int64]int, len(users))
		for i := 0; i < len(users); i++ { // 记录 用户 id 在 数组的位置
			idx[users[i].Id] = i
		}
	}

	res := make([]*handlers.Video, len(videos))
	for i := 0; i < len(videos); i++ {

		res[i] = &handlers.Video{
			Id:            videos[i].Id,
			PlayUrl:       videos[i].PlayUrl,
			CoverUrl:      videos[i].CoverUrl,
			FavoriteCount: videos[i].FavoriteCount,
			CommentCount:  videos[i].CommentCount,
			IsFavorite:    videos[i].IsFavorite,
		}

		if isOne {
			res[i].Author = one
			continue
		}

		if j, ok := idx[videos[i].Author.Id]; ok {
			res[i].Author = users[j]
		} else {
			klog.Infof("[Pack] 无法找到对应 id; users: %d Author: %d", j, videos[i].Author)
		}
	}
	return res

}
