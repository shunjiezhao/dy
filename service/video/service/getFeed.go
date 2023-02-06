package service

import (
	"context"
	"first/kitex_gen/video"
	"first/service/video/model/db"
	"first/service/video/pack"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

type FeedsService struct {
	ctx context.Context
	//point time.Time在这个时间点之前的 limit 个 视频url
	//limit int
}

// NewFeedsService new CreateNoteService
func NewFeedsService(ctx context.Context) *FeedsService {
	return &FeedsService{
		ctx: ctx,
	}
}

// FeedsItem 返回视频列表, 前提是没有用户id,也就是不需要查询是否喜欢
func (s *FeedsService) FeedsItem(req *video.GetVideoListRequest) ([]*video.Video, error) {

	videos, err := db.GetVideosAfterTime(s.ctx, req.TimeStamp, 30)
	if err != nil {
		klog.Infof("[DB]: 得到 Feeds 流失败; err: %v", err.Error())
		return nil, err
	}
	return pack.Videos(videos), nil
}

// GetUserPublish create note info
func (s *FeedsService) GetUserPublish(req *video.GetVideoListRequest) ([]*video.Video, error) {
	videos, err := db.GetUserPublish(db.VideoDb, s.ctx, req.Author)
	if err != nil {
		klog.Infof("[DB]: 得到 Feeds 流失败; err: %v", err.Error())
		return nil, err
	}
	return pack.Videos(videos), nil
}

// LoginUserFeedsItem 返回给登陆用户视频列表, 需要查询是否喜欢
func (s *FeedsService) LoginUserFeedsItem(req *video.GetVideoListRequest) ([]*video.Video, error) {
	if req.TimeStamp == 0 {
		req.TimeStamp = time.Now().Unix()
	}

	videos, err := db.LoginUserFeedsItem(s.ctx, req.TimeStamp, req.Uuid)
	if err != nil {
		klog.Infof("[DB]: 登陆用户得到 Feeds 流失败; err: %v", err.Error())
		return nil, err
	}
	return pack.LoginFeeds(videos), nil
}
