package service

import (
	"context"
	"first/kitex_gen/video"
	"first/pkg/util"
	"first/service/video/model/db"
)

type CreateVideoItemService struct {
	ctx context.Context
}

// NewCreateVideoItemService new CreateNoteService
func NewCreateVideoItemService(ctx context.Context) *CreateVideoItemService {
	return &CreateVideoItemService{ctx: ctx}
}

// CreateVideoItem create note info
func (s *CreateVideoItemService) CreateVideoItem(req *video.PublishListRequest) error {
	dVideo := &db.Video{
		Id:         util.NextVal(),
		AuthorUuid: req.Author,
		Title:      req.Title,
		PlayUrl:    req.PlayUrl,
		CoverUrl:   req.CoverUrl,
	}

	return db.CreateVideoItem(db.VideoDb, s.ctx, dVideo)
}
