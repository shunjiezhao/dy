package service

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/util"
	userDB "first/service/user/model/db"
)

type CommentService struct {
	ctx context.Context
}

// NewCommentService new CreateNoteService
func NewCommentService(ctx context.Context) *CommentService {
	return &CommentService{ctx: ctx}
}

// CreateComment 新建评论
func (s *CommentService) CreateComment(req *userPb.ActionCommentRequest) error {
	return userDB.CreateComment(s.ctx, &userDB.Comment{
		Id:      util.NextVal(),
		Uuid:    req.Uuid,
		VideoId: req.VideoId,
		Content: req.CommentText,
	})
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(req *userPb.ActionCommentRequest) error {
	return userDB.DeleteComment(s.ctx, req.CommentId)
}
