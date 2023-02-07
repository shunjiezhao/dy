package service

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/util"
	userDB "first/service/user/model/db"
	"first/service/user/pack"
)

type CommentService struct {
	ctx context.Context
}

// NewCommentService new CreateNoteService
func NewCommentService(ctx context.Context) *CommentService {
	return &CommentService{ctx: ctx}
}

// CreateComment 新建评论
func (s *CommentService) CreateComment(req *userPb.ActionCommentRequest) (*userPb.Comment, error) {
	comment, err := userDB.CreateComment(s.ctx, &userDB.Comment{
		Id:      util.NextVal(),
		Uuid:    req.Uuid,
		VideoId: req.VideoId,
		Content: req.CommentText,
	})
	if err != nil {
		return nil, err
	}

	return pack.Comment(comment), nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(req *userPb.ActionCommentRequest) error {
	return userDB.DeleteComment(s.ctx, req.CommentId)
}
