package service

import (
	userPb "first/kitex_gen/user"
	userDB "first/service/user/model/db"
	"first/service/user/pack"
	"time"
)

// GetComment 新建评论
func (s *CommentService) GetComment(req *userPb.GetCommentRequest) ([]*userPb.Comment, error) {
	comments, err := userDB.GetCommentAfterTime(s.ctx, req.VideoId, time.Now().Unix(), 10)
	if err != nil {
		return nil, err
	}
	return pack.Comments(comments), nil
}
