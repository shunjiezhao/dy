package service

import (
	"context"
	"first/kitex_gen/video"
	"first/service/video/model/db"
	"first/service/video/pack"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

type LikeService struct {
	ctx context.Context
}

// NewLikeService new CreateNoteService
func NewLikeService(ctx context.Context) *LikeService {
	return &LikeService{
		ctx: ctx,
	}
}

// LikeVideo 喜欢视频
func (s *LikeService) LikeVideo(req *video.LikeVideoRequest) error {
	param := &db.FavouriteVideo{
		Uuid:    req.Uuid,
		VideoId: req.VideoId,
		IsLike:  true,
	}
	err := db.CreateFavVideo(s.ctx, param)
	if err != nil {
		klog.Infof("[Video-DB]: 喜欢操作失败; err: %v", err.Error())
		return err
	}
	return nil
}

// UnLikeVideo 不喜欢视频
func (s *LikeService) UnLikeVideo(req *video.LikeVideoRequest) error {
	param := &db.FavouriteVideo{
		Uuid:    req.Uuid,
		VideoId: req.VideoId,
	}
	err := db.DeleteFavVideo(s.ctx, param)

	if err != nil {
		klog.Infof("[Video-DB]: 不喜欢操作失败; err: %v", err.Error())
		return err
	}
	return nil
}

// LikesItem 返回用户喜欢的视频列表
func (s *LikeService) LikesItem(req *video.GetVideoListRequest) ([]*video.Video, error) {
	videos, err := db.GetFavVideoAfterTime(s.ctx, req.Uuid, time.Now().Unix(), 20)
	if err != nil {
		klog.Infof("[Video-DB]: 得到  喜欢列表 失败; err: %v", err.Error())
		return nil, err
	}

	return pack.FavVideos(videos), nil
}
