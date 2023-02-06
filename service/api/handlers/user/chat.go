package user

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func (s *Service) SendMsg() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err       error
			param     ChatActionRequest //http 请求参数
			ctx       context.Context   = c.Request.Context()
			curUserId int64
			msg       handlers.Message
		)

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetToUserId() == 0 || len(param.GetToken()) == 0 {
			err = c.ShouldBind(&param)
		}
		if err != nil {
			klog.Errorf("[发送消息]: 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		curUserId = getTokenUserId(c)
		// 不能发送消息给自己
		if param.GetToUserId() == curUserId {
			handlers.SendResponse(c, errno.OpSelfErr)
			goto errHandler

		}

		//TODO: msg id
		klog.Infof("[%d->%d]: %s", curUserId, param.GetToUserId(), param.Content)
		msg.ToUserId = param.ToUserId
		msg.FromUserId = handlers.FromUserId{UserId: curUserId}
		msg.Content = param.Content
		msg.CreateTime = time.Now().Format(constants.TimeFormatS)

		//TODO: 优化 考虑如果消息发送失败
		err = s.chatSrv.Save(ctx, &msg)
		if err != nil {
			klog.Errorf("[保存消息]: 失败 %v", err.Error())
			handlers.SendResponse(c, err)
			goto errHandler

		}

		SendChatResponse(c, &msg)
		return
	errHandler:
		c.Abort()
	}
}

func (s *Service) GetChatList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err       error
			param     ChatListRequest //http 请求参数
			ctx       context.Context = c.Request.Context()
			curUserId int64
			msgs      []*handlers.Message
		)

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetToUserId() == 0 || len(param.GetToken()) == 0 {
			err = c.ShouldBind(&param)
		}

		if err != nil {
			klog.Errorf("[获取消息]: 绑定参数失败 %v", err.Error())
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		curUserId = getTokenUserId(c)
		// 查看自己消息
		if param.GetToUserId() == curUserId {
			handlers.SendResponse(c, errno.OpSelfErr)
			goto errHandler

		}

		msgs, err = s.chatSrv.GetList(ctx, handlers.FromUserId{UserId: curUserId}, param.ToUserId)
		if err != nil {
			klog.Errorf("[获取消息] 获取消息失败 %v", err.Error())
			handlers.SendResponse(c, err)
			goto errHandler

		}
		klog.Infof("[获取消息]: [%d->%d]聊天 信息", curUserId, param.GetToUserId())
		log.Printf("%+v", msgs[0])
		log.Printf("%#v", msgs[0])

		GetChatListResponse(c, msgs)
		return
	errHandler:
		c.Abort()
	}
}
