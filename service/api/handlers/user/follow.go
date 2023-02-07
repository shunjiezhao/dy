package user

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
)

var shouldGetOther = false

// 返回是否有权限得到其他人的关注信息
func getOther(userId, otherId int64) bool {
	if userId == otherId {
		return true
	}
	return shouldGetOther
}

//GetFriendList 获取好友列表
func (s *Service) GetFriendList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err       error
			param     common.FriendListRequest //http 请求参数
			curUserId int64                    //当前用户的 userid
			list      []*handlers.Message      // 返回的粉丝列表
			ctx       context.Context          = c.Request.Context()
		)
		curUserId = getTokenUserId(c)
		if curUserId == -1 {
			klog.Error("获取 user_id 出错")
			goto errHandler
		}

		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetUserId() == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		if err != nil {
			klog.Errorf("[获取好友列表]: 获取 参数 %+v", param)
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		// rpc 调用
		list, err = s.chatSrv.GetFriendChatList(ctx, handlers.FromUserId{UserId: curUserId})
		if err != nil {
			klog.Errorf("[获取好友列表]: 调用 rpc 失败%v", err)
			handlers.SendResponse(c, err)
			goto errHandler
		}

		common.GetChatListResponse(c, list)
		return
	errHandler:
		c.Abort()
	}
}
func (s *Service) GetFollowerList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err       error
			param     common.GetUserFollowerListRequest //http 请求参数
			curUserId int64                             //当前用户的 userid
			list      []*handlers.User                  // 返回的粉丝列表
			ctx       context.Context                   = c.Request.Context()
		)
		curUserId = getTokenUserId(c)
		if curUserId == -1 {
			klog.Error("获取 user_id 出错")
			goto errHandler
		}
		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetUserId() == 0 || len(param.GetToken()) == 0 {
			err = c.ShouldBind(&param)
		}
		if err != nil || !getOther(curUserId, param.GetUserId()) {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler

		}

		klog.Infof("[粉丝列表]: 获取 参数 %+v", param)

		// rpc 调用

		list, err = s.rpc.GetFollowerList(ctx, &param)
		if err != nil {
			klog.Errorf("[粉丝列表]: 调用 rpc 失败%v", err)
			handlers.SendResponse(c, err)
			goto errHandler

		}
		common.SendUserListResponse(c, list)
		return
	errHandler:
		c.Abort()
	}
}
func (s *Service) GetFollowList() func(c *gin.Context) {

	return func(c *gin.Context) {
		var (
			err       error
			param     common.GetUserFollowListRequest //http 请求参数
			curUserId int64                           //当前用户的 userid
			list      []*handlers.User                // 返回的关注列表
			ctx       context.Context                 = c.Request.Context()
		)
		curUserId = getTokenUserId(c)
		if curUserId == -1 {
			klog.Error("获取 user_id 出错")
			goto errHandler

		}
		// token 检验成功 开始  绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetUserId() == 0 || len(param.GetToken()) == 0 {
			err = c.ShouldBind(&param) // bind form
		}

		if err != nil || !getOther(curUserId, param.GetUserId()) {
			klog.Errorf("[关注列表]: 只能获取自己的 %d", curUserId)
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		klog.Infof("[关注列表]: 获取 参数 %+v", param)

		// rpc

		list, err = s.rpc.GetFollowList(ctx, &param)
		if err != nil {
			klog.Errorf("[关注列表]: 调用 rpc 失败%v", err)
			handlers.SendResponse(c, err)
			goto errHandler

		}
		common.SendUserListResponse(c, list)
		return
	errHandler:
		c.Abort()
	}
}
func (s *Service) Follow() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			param common.ActionRequest
			err   error
			ctx   context.Context = c.Request.Context()
		)
		curUserId := getTokenUserId(c)
		if curUserId == -1 {
			klog.Error("获取 user_id 出错")
			goto errHandler
		}
		// token 检验成功
		// 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.ToUserId == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		// 当前用户不能关注自己
		if err != nil || curUserId == param.ToUserId {
			handlers.SendResponse(c, errno.OpSelfErr)
			goto errHandler
		}
		klog.Infof("[关注操作]: 获取 参数 %+v", param)

		param.FromUserId = curUserId
		// 发送绑定请求
		err = s.rpc.ActionFollow(ctx, &param)
		if err != nil { // remote  network error
			klog.Errorf("[关注操作]: 调用 rpc 失败%v", err)
			handlers.SendResponse(c, err)
			goto errHandler
		}

		handlers.SendResponse(c, errno.Success)
		return
	errHandler:
		c.Abort()
	}
}

func getTokenUserId(c *gin.Context) int64 {
	claim, ok := c.Get(constants.IdentityKey)
	if !ok {
		return -1
	}

	var curUserId int64
	tmp, ok := claim.(float64)
	if ok {
		curUserId = int64(tmp)
	} else {
		curUserId = claim.(int64)
	}
	return curUserId
}
