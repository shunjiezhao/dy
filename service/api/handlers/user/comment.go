package user

import (
	"context"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

func (s *Service) GetCommentList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err      error
			param    common.CommentListRequest //http 请求参数
			ctx      context.Context           = c.Request.Context()
			comments []*handlers.Comment
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

		comments, err = s.rpc.GetComment(ctx, &param)
		if err != nil {
			klog.Errorf("[获取评论]: 调用[用户服务] 获取评论失败 %v", err.Error())
			handlers.SendResponse(c, errno.RemoteErr)
			goto errHandler
		}

		common.SendCommentListResponse(c, comments)
		return
	errHandler:
		c.Abort()
	}
}

func (s *Service) ActionComment() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err     error
			param   common.CommentActionRequest //http 请求参数
			ctx     context.Context             = c.Request.Context()
			comment *handlers.Comment
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
		param.UserId = getTokenUserId(c)

		klog.Infof("[%d->%d]: %s评论", param.UserId, param.VideoId, param.CommentActionType)

		// rpc 调用

		comment, err = s.rpc.ActionComment(ctx, &param)
		if err != nil {
			klog.Errorf("[评论操作] 调用[用户服务] 获取评论失败 %v", err.Error())
			handlers.SendResponse(c, err)
			goto errHandler
		}

		common.SendCommentResponse(c, comment)
		return
	errHandler:
		c.Abort()
	}
}
