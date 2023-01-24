package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/rpc"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"time"
)

// Register 注册用户
// @tokenGenerator: 生成 token
// @RpcRegister: 调用 rpc 返回 userid
func Register(tokenGenerator func(data interface{}) (string, time.Time, error), RpcRegister func(ctx context.Context, req *userPb.RegisterRequest) (int64, error)) func(ctx context.Context,
	c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {

		var param RegisterRequest
		err := c.Bind(&param)
		// 参数校验
		if err != nil || len(param.UserName) == 0 || len(param.PassWord) == 0 {
			handlers.SendResponse(c, err)
			return
		}

		req := userPb.RegisterRequest{
			UserName: param.UserName,
			PassWord: param.PassWord,
		}
		if RpcRegister == nil {
			RpcRegister = rpc.Register
		}
		userId, err := RpcRegister(ctx, &req) // 方便mock
		if err != nil {
			handlers.SendResponse(c, err)
			return
		}
		token, _, err := tokenGenerator(userId)
		if err != nil {
			handlers.SendResponse(c, err)
			return
		}
		SendRegisterResponse(c, userId, token)
	}
}
func Login() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		var loginVar LoginRequest
		err := c.Bind(&loginVar)
		if err != nil || len(loginVar.UserName) == 0 || len(loginVar.PassWord) == 0 {
			handlers.SendResponse(c, errno.ParamErr)
			c.Abort()
			return
		}
		req := &userPb.CheckUserRequest{
			UserName: loginVar.UserName,
			PassWord: loginVar.PassWord,
		}
		Uuid, err := rpc.CheckUser(ctx, req)
		if err != nil {
			handlers.SendResponse(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		}
		c.Set(constants.IdentityKey, Uuid)
	}
}

// 得到里面的 user_id
func checkToken(jwt *jwt.HertzJWTMiddleware, ctx context.Context, c *app.RequestContext) (int64, error) {
	// 检查token
	fromJWT, err := jwt.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return 0, errno.ParamErr
	}
	userId, ok := fromJWT[constants.IdentityKey]
	if !ok { // token 中没有 user_id 不合法
		return 0, errno.AuthorizationFailedErr
	}
	// 这里并不知道为什么 传递的是 int64 但是 存放的是 float64
	id, ok := userId.(float64) // 断言失败
	if !ok {
		return 0, errno.ParamErr
	}
	return int64(id), nil
}
func GetInfo(jwt *jwt.HertzJWTMiddleware) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查参数
		var param GetInfoRequest
		err := c.Bind(&param)
		if err != nil || param.UserId.UserId <= int64(0) {
			handlers.SendResponse(c, errno.ParamErr)
			return
		}
		_, err = checkToken(jwt, ctx, c)
		if err != nil {
			handlers.SendResponse(c, err)
			return
		}
		// 发送查询请求
		req := &userPb.GetUserRequest{Id: param.UserId.UserId}
		user, err := rpc.GetUserInfo(ctx, req)
		if err != nil {
			handlers.SendResponse(c, err)
			return
		}
		SendGetInfoResponse(c, handlers.PackUser(user))
	}
}
