package service

import (
	"context"
	"first/kitex_gen/video"
	"first/pkg/util"
	"first/service/video/model/db"
	"github.com/cloudwego/kitex/pkg/klog"
)

type VideoItemService struct {
	ctx context.Context
}

// NewVideoItemService new CreateNoteService
func NewVideoItemService(ctx context.Context) *VideoItemService {
	return &VideoItemService{ctx: ctx}
}

// CreateVideoItem 新建视频信息
func (s *VideoItemService) CreateVideoItem(req *video.PublishListRequest) error {
	dVideo := &db.Video{
		Id:         util.NextVal(),
		AuthorUuid: req.Author,
		Title:      req.Title,
		PlayUrl:    req.PlayUrl,
		CoverUrl:   req.CoverUrl,
	}
	klog.Infof("保存信息%#v", req)
	return db.CreateVideoItem(db.VideoDb, s.ctx, dVideo)
}

// IncrCommentCount 更新视频评论数
func (s *VideoItemService) IncrCommentCount(req *video.IncrCommentRequest) error {

	return db.IncrVideoCommentCount(db.VideoDb, s.ctx, req.VideoId, req.Add)
}
