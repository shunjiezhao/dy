package service

import (
	"context"
	"first/kitex_gen/user"
	"first/service/user/model/db"
	"log"
)

type FollowUserService struct {
	ctx context.Context
}

// NewFollowUserService new CreateNoteService
func NewFollowUserService(ctx context.Context) *FollowUserService {
	return &FollowUserService{ctx: ctx}
}

// FollowUser create note info
func (s *FollowUserService) FollowUser(req *user.FollowRequest) (bool, error) {
	log.Println("FollowUser: rpc 响应开始调用")
	return db.FollowUser(s.ctx, req.ToUserId, req.FromUserId)
}

type UnFollowUserService struct {
	ctx context.Context
}

// NewUnFollowUserService new CreateNoteService
func NewUnFollowUserService(ctx context.Context) *UnFollowUserService {
	return &UnFollowUserService{ctx: ctx}
}

// UnFollowUser create note info
func (s *UnFollowUserService) UnFollowUser(req *user.FollowRequest) (bool, error) {
	log.Println("UnFollowUser: rpc 响应开始调用")
	return db.UnFollowUser(s.ctx, req.ToUserId, req.FromUserId)
}
