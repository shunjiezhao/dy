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

func GetFollowerList(parse func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		curUserId := getTokenUserId(ctx, c, parse)
		if curUserId == -1 {
			return
		}
		// token 检验成功
		var param GetUserFollowerListRequest
		// 绑定参数
		err := c.Bind(&param)
		//TODO: 这里是不能获取别人的关注列表, 但是 是否可以 真正获取别人的关注列表呢?
		if err != nil || curUserId != param.UserId.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			c.Abort()
			return
		}
		// rpc
		req := &userPb.GetFollowerListRequest{
			Id: param.UserId.UserId,
		}
		list, err := rpc.GetFollowerList(ctx, req)
		if err != nil {
			handlers.SendResponse(c, err)
			c.Abort()
			return
		}
		SendUserListResponse(c, handlers.PackUsers(list))
	}
}
func GetFollowList(parse func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		curUserId := getTokenUserId(ctx, c, parse)
		if curUserId == -1 {
			return
		}
		// token 检验成功
		var param GetUserFollowListRequest
		// 绑定参数
		err := c.Bind(&param)
		//TODO: 这里是不能获取别人的关注列表, 但是 是否可以 真正获取别人的关注列表呢?
		if err != nil || curUserId != param.UserId.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			c.Abort()
			return
		}
		// rpc
		req := &userPb.GetFollowListRequest{
			Id: param.UserId.UserId,
		}
		list, err := rpc.GetFollowList(ctx, req)
		if err != nil {
			handlers.SendResponse(c, err)
			c.Abort()
			return
		}
		SendUserListResponse(c, handlers.PackUsers(list))
	}
}
func Follow(parse func(ctx context.Context, c *app.RequestContext) (jwt2.MapClaims, error)) func(ctx context.Context,
	c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		curUserId := getTokenUserId(ctx, c, parse)
		if curUserId == -1 {
			return
		}
		// token 检验成功
		// 判断 token 中的 Id 是否是自己
		var param ActionRequest
		// 绑定参数
		err := c.Bind(&param)
		// 当前用户不能关注自己
		if err != nil || curUserId == param.UserId {
			handlers.SendResponse(c, errno.ParamErr)
			c.Abort()
			return
		}
		// 发送绑定请求
		var req = &userPb.FollowRequest{
			FromUserId: curUserId,
			ToUserId:   param.UserId,
		}
		if param.IsFollow() {
			err = rpc.FollowUser(ctx, req)
		} else {
			err = rpc.UnFollowUser(ctx, req)
		}
		if err != nil {
			handlers.SendResponse(c, err)
			c.Abort()
			return
		}
		handlers.SendResponse(c, errno.Success)
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
