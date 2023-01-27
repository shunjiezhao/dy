package follow

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/rpc"
	"github.com/cloudwego/hertz/pkg/app"
	jwt2 "github.com/hertz-contrib/jwt"
)

type Service struct {
	parseToken func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)
}

func New(parse func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)) *Service {
	return &Service{parseToken: parse}
}

//TODO: 1. 不能获取别人的关注/粉丝列表, 但是可以获取别人的关注/粉丝列表嘛?

func (s *Service) GetFollowerList() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		var (
			err       error
			req       *userPb.GetFollowerListRequest // rpc 调用参数
			param     GetUserFollowerListRequest     //http 请求参数
			curUserId int64                          //当前用户的 userid
			list      []*userPb.User                 // 返回的粉丝列表
		)
		curUserId = getTokenUserId(ctx, c, s.parseToken)
		if curUserId == -1 {
			goto errHandler
		}
		// token 检验成功 开始 绑定参数
		err = c.Bind(&param)
		//TODO: 1
		if err != nil || curUserId != param.UserId.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		// rpc 调用
		req = &userPb.GetFollowerListRequest{
			Id: param.UserId.UserId,
		}
		list, err = rpc.GetFollowerList(ctx, req)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler
		}
		SendUserListResponse(c, handlers.PackUsers(list))
	errHandler:
		c.Abort()
		return
	}
}
func (s *Service) GetFollowList() func(ctx context.Context, c *app.RequestContext) {

	return func(ctx context.Context, c *app.RequestContext) {
		var (
			err       error
			req       *userPb.GetFollowListRequest // rpc 调用参数
			param     GetUserFollowListRequest     //http 请求参数
			curUserId int64                        //当前用户的 userid
			list      []*userPb.User               // 返回的关注列表
		)
		curUserId = getTokenUserId(ctx, c, s.parseToken)
		if curUserId == -1 {
			return
		}
		// token 检验成功 开始  绑定参数
		err = c.Bind(&param)
		//TODO: 1
		if err != nil || curUserId != param.UserId.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}
		// rpc
		req = &userPb.GetFollowListRequest{
			Id: param.UserId.UserId,
		}
		list, err = rpc.GetFollowList(ctx, req)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler
		}
		SendUserListResponse(c, handlers.PackUsers(list))
	errHandler:
		c.Abort()
		return
	}
}
func (s *Service) Follow() func(ctx context.Context,
	c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		var (
			param ActionRequest
			req   *userPb.FollowRequest
			err   error
		)
		curUserId := getTokenUserId(ctx, c, s.parseToken)
		if curUserId == -1 {
			goto errHandler
		}
		// token 检验成功
		// 绑定参数
		err = c.Bind(&param)
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
			err = rpc.FollowUser(ctx, req)
		} else {
			err = rpc.UnFollowUser(ctx, req)
		}
		if err != nil { // remote  network error
			handlers.SendResponse(c, err)
			goto errHandler
		}
		handlers.SendResponse(c, errno.Success)
	errHandler:
		c.Abort()
		return
	}
}

func getTokenUserId(ctx context.Context, c *app.RequestContext, parse func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)) int64 {
	claim, err := parse(ctx, c)
	if err != nil {
		handlers.SendResponse(c, errno.AuthorizationFailedErr)
		c.Abort()
		return -1
	}
	var curUserId int64
	tmp, ok := claim[constants.IdentityKey].(float64)
	if ok {
		curUserId = int64(tmp)
	} else {
		curUserId = claim[constants.IdentityKey].(int64)
	}
	return curUserId
}
