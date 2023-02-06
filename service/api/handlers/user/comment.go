package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) GetCommentList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err      error
			req      *userPb.GetCommentRequest // rpc 调用参数
			param    CommentListRequest        //http 请求参数
			ctx      context.Context           = c.Request.Context()
			comments []*userPb.Comment
		)

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.VideoId == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		if err != nil {
			klog.Errorf("[获取评论] 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		// rpc 调用
		req = &userPb.GetCommentRequest{
			VideoId: param.VideoId,
		}

		comments, err = s.rpc.GetComment(ctx, req)
		if err != nil {
			klog.Errorf("[获取评论]: 调用[用户服务] 获取评论失败 %v", err.Error())
			handlers.SendResponse(c, errno.RemoteErr)
			goto errHandler
		}

		SendCommentListResponse(c, handlers.PackComments(comments))
		return
	errHandler:
		c.Abort()
	}
}

func (s *Service) ActionComment() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err     error
			req     *userPb.ActionCommentRequest // rpc 调用参数
			param   CommentActionRequest         //http 请求参数
			ctx     context.Context              = c.Request.Context()
			comment *userPb.Comment
		)

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.VideoId == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		if err != nil {
			klog.Errorf("[评论操作]: 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		klog.Infof("[%d->%d]: %s评论", handlers.GetTokenUserId(c), param.VideoId, param.CommentActionType)

		// rpc 调用
		req = &userPb.ActionCommentRequest{
			Uuid:        handlers.GetTokenUserId(c),
			VideoId:     param.VideoId,
			ActionType:  &userPb.MessageActionType{Op: int32(param.CommentActionType)},
			CommentText: param.CommentText,
			CommentId:   param.CommentId,
		}
		comment, err = s.rpc.ActionComment(ctx, req)
		if err != nil {
			klog.Errorf("[评论操作] 调用[用户服务] 获取评论失败 %v", err.Error())
			handlers.SendResponse(c, err)
			goto errHandler
		}

		SendCommentResponse(c, handlers.PackComment(comment))
		return
	errHandler:
		c.Abort()
	}
}
