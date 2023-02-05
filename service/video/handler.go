package main

import (
	"context"
	video "first/kitex_gen/video"
	"first/pkg/errno"
	"first/service/video/pack"
	"first/service/video/service"
	"github.com/cloudwego/kitex/pkg/klog"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

func (s *VideoServiceImpl) LikeVideo(ctx context.Context, req *video.LikeVideoRequest) (resp *video.LikeVideoResponse, err error) {
	resp = new(video.LikeVideoResponse)
	if req.VideoId == 0 || req.ActionType == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	if req.ActionType.ActionType == 1 {
		err = service.NewLikeService(ctx).LikeVideo(req)
	} else {
		err = service.NewLikeService(ctx).UnLikeVideo(req) // 不喜欢
	}
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.LikeOpErr)
		return resp, nil

	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return resp, nil

}

func (s *VideoServiceImpl) GetLikeVideo(ctx context.Context, req *video.GetVideoListRequest) (resp *video.GetVideoListResponse, err error) {
	resp = new(video.GetVideoListResponse)
	if req.Author == 0 || req.GetAuthor_ == false {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}

	resp.VideoList, err = service.NewLikeService(ctx).LikesItem(req) // 获取喜欢列表
	klog.Infof("get video list :%#v", resp.VideoList)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.GetVideoErr)
		return resp, nil
	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// Upload implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Upload(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	resp := new(video.PublishListResponse)
	err := service.NewCreateVideoItemService(ctx).CreateVideoItem(req) // 创建 video item
	if err != nil {
		klog.Errorf("save video item error: %v", err.Error())
		resp.Resp = pack.BuildBaseResp(errno.PublishVideoErr)
	} else {
		resp.Resp = pack.BuildBaseResp(errno.Success)
	}
	return resp, nil
}

// GetVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideoList(ctx context.Context, req *video.GetVideoListRequest) (*video.GetVideoListResponse, error) {
	resp := new(video.GetVideoListResponse)
	var (
		err error
	)
	if req.GetAuthor_ {
		resp.VideoList, err = service.NewFeedsService(ctx).GetUserPublish(req) // 用户发布的
	} else if req.Uuid == 0 {
		resp.VideoList, err = service.NewFeedsService(ctx).FeedsItem(req) // 未登录用户
	} else {
		resp.VideoList, err = service.NewFeedsService(ctx).LoginUserFeedsItem(req) // 登陆用户获取, 需要 传递是否喜欢
	}
	klog.Infof("get video list :%#v", resp.VideoList)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.GetVideoErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}
