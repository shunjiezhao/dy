package user

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/mq"
	"first/pkg/util"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"github.com/u2takey/go-utils/json"
	"time"
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
		if err != nil || param.VideoId == 0 || len(param.GetToken()) == 0 {
			klog.Errorf("[获取评论] 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		klog.Errorf("[获取评论] 参数 %+v", param)

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
			data    []byte
		)

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil {
			err = c.ShouldBind(&param)
		}
		if err != nil || len(param.GetToken()) == 0 {
			klog.Errorf("[评论操作]: 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		param.UserId = getTokenUserId(c)

		if param.CommentActionType.IsAdd() {
			param.CommentId = util.NextVal()
		}
		// 发送消息队列
		data, err = param2Info(&param)
		if err != nil {
			klog.Errorf("[评论操作] json 化失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler

		}

		err = s.publisher[constants.UActionCommentKey][mq.UGetActionCommentIdx(param.VideoId)].Publish(ctx, data)
		if err != nil {
			klog.Infof("[评论操作]发送失败")
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		err = s.publisher[constants.VActionVideoComCountKey][mq.VGetActionVideoComCountIdx(param.VideoId)].Publish(ctx, data)
		if err != nil {
			klog.Infof("[评论操作]发送失败")
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler

		}

		// rpc 调用

		comment = &handlers.Comment{
			Id:         param.CommentId,
			Content:    param.CommentText,
			CreateDate: time.Now().Format(constants.TimeFormatS),
		}
		klog.Infof("[%d->%d]: %s评论", param.UserId, param.VideoId, param.CommentActionType)
		klog.Infof("[%d->%d]: %s评论", param.UserId, param.VideoId, comment)

		comment.User = &handlers.User{
			Id:   param.UserId,
			Name: c.MustGet(constants.UserNameKey).(string),
		}
		common.SendCommentResponse(c, comment)
		return
	errHandler:
		c.Abort()
	}
}
func param2Info(param *common.CommentActionRequest) ([]byte, error) {
	info := mq.ActionCommentInfo{
		Uuid:        param.UserId,
		VideoId:     param.VideoId,
		ActionType:  int32(param.CommentActionType),
		CommentText: param.CommentText,
		CommentId:   param.CommentId,
	}
	return json.Marshal(info)
}
