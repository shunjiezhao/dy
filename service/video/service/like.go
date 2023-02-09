package service

import (
	"context"
	"first/kitex_gen/video"
	"first/pkg/redis"
	video2 "first/pkg/redis/video"
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
	go func() {
		time.Sleep(time.Millisecond * 500) // 500ms
		if err := video2.DelFavVideoList(redis.GetRedis(), s.ctx, req.Uuid); err != nil {
			klog.Errorf("删除失败")
		}
	}()

	// 双删
	video2.DelFavVideoList(redis.GetRedis(), s.ctx, req.Uuid)

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
	go func() {
		time.Sleep(time.Millisecond * 500) // 500ms
		if err := video2.DelFavVideoList(redis.GetRedis(), s.ctx, req.Uuid); err != nil {
			klog.Errorf("删除失败")
		}
	}()

	// 双删
	video2.DelFavVideoList(redis.GetRedis(), s.ctx, req.Uuid)
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
	// Redis 命中
	list, err := video2.GetFavVideoList(redis.GetRedis(), s.ctx, req.Uuid)
	if err == nil {
		return list, err
	}
	videos, err := db.GetFavVideoAfterTime(s.ctx, req.Uuid, time.Now().Unix(), 20)
	if err != nil {
		klog.Infof("[Video-DB]: 得到  喜欢列表 失败; err: %v", err.Error())
		return nil, err
	}
	ans := pack.FavVideos(videos)
	video2.WriteFavVideo(redis.GetRedis(), s.ctx, ans, req.Uuid) // 写入
	return ans, nil
}
