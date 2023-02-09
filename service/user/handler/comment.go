package handler

import (
	"context"
	"first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/pack"
	comment "first/service/user/service/comment"
	"github.com/cloudwego/kitex/pkg/klog"
)

// isAdd 是否是发布评论
func isAdd(i int32) bool {
	if i == 1 {
		return true
	}
	return false
}
func (s *UserServiceImpl) ActionComment(ctx context.Context, req *user.ActionCommentRequest) (resp *user.
	ActionCommentResponse, err error) {
	resp = new(user.ActionCommentResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return

	}
	if isAdd(req.ActionType) { // 创建
		resp.Comment, err = comment.NewCommentService(ctx).CreateComment(req)
	} else {
		err = comment.NewCommentService(ctx).DeleteComment(req)
	}

	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	klog.Infof("操作成功 %+v", req)
	resp.Resp = pack.BuildBaseResp(errno.Success)

	return
}

func (s *UserServiceImpl) GetComment(ctx context.Context, req *user.GetCommentRequest) (resp *user.GetCommentResponse, err error) {
	resp = new(user.GetCommentResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		klog.Infof("[GetComment]: 参数有误")
		return resp, nil

	}
	resp.Comment, err = comment.NewCommentService(ctx).GetComment(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.UserAlreadyExistErr)
		return resp, nil
	}
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
}
