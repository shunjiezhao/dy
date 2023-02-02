package user

import (
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/gin-gonic/gin"
)

//TODO: 1. 不能获取别人的关注/粉丝列表, 但是可以获取别人的关注/粉丝列表嘛?
//  	2. 和 user srv 一样 将所有的 rpc 调用抽离开来, 方便mock

func (s *Service) GetFollowerList() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err       error
			req       *userPb.GetFollowerListRequest // rpc 调用参数
			param     GetUserFollowerListRequest     //http 请求参数
			curUserId int64                          //当前用户的 userid
			list      []*userPb.User                 // 返回的粉丝列表
		)
		curUserId = getTokenUserId(c)
		if curUserId == -1 {
			goto errHandler
		}
		// token 检验成功 开始 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetUserId() == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		//TODO: 1
		if err != nil || curUserId != param.GetUserId() {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		// rpc 调用
		req = &userPb.GetFollowerListRequest{
			Id: param.GetUserId(),
		}
		list, err = s.rpc.GetFollowerList(c, req)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler
		}
		SendUserListResponse(c, handlers.PackUsers(list))
		return
	errHandler:
		c.Abort()
	}
}
func (s *Service) GetFollowList() func(c *gin.Context) {

	return func(c *gin.Context) {
		var (
			err       error
			req       *userPb.GetFollowListRequest // rpc 调用参数
			param     GetUserFollowListRequest     //http 请求参数
			curUserId int64                        //当前用户的 userid
			list      []*userPb.User               // 返回的关注列表
		)
		curUserId = getTokenUserId(c)
		if curUserId == -1 {
			return
		}
		// token 检验成功 开始  绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.GetUserId() == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param) // bind form
		}

		//TODO: 1
		if err != nil || curUserId != param.GetUserId() {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		// rpc
		req = &userPb.GetFollowListRequest{
			Id: param.GetUserId(),
		}
		list, err = s.rpc.GetFollowList(c, req)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler
		}
		SendUserListResponse(c, handlers.PackUsers(list))
		return
	errHandler:
		c.Abort()
	}
}
func (s *Service) Follow() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			param ActionRequest
			req   *userPb.FollowRequest
			err   error
		)
		curUserId := getTokenUserId(c)
		if curUserId == -1 {
			goto errHandler
		}
		// token 检验成功
		// 绑定参数
		err = c.ShouldBindQuery(&param)
		if err != nil || param.UserId == 0 || param.GetToken() == "" {
			err = c.ShouldBind(&param)
		}
		// 当前用户不能关注自己
		if err != nil || curUserId == param.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		// 发送绑定请求
		req = &userPb.FollowRequest{
			FromUserId: curUserId,
			ToUserId:   param.UserId,
		}
		if param.IsFollow() {
			err = s.rpc.FollowUser(c, req)
		} else {
			err = s.rpc.UnFollowUser(c, req)
		}
		if err != nil { // remote  network error
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
	claim := c.MustGet(constants.IdentityKey)

	var curUserId int64
	tmp, ok := claim.(float64)
	if ok {
		curUserId = int64(tmp)
	} else {
		curUserId = claim.(int64)
	}
	return curUserId
}
