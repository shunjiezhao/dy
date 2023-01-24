package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/rpc"
	"github.com/cloudwego/hertz/pkg/app"
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
			handlers.SendResponse(c, err, nil)
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
			handlers.SendResponse(c, err, nil)
			return
		}
		//TODO: jwt 生成token
		token, _, err := tokenGenerator(userId)
		if err != nil {
			handlers.SendResponse(c, err, nil)
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
			handlers.SendResponse(c, errno.ParamErr, nil)
			c.Abort()
			return
		}
		req := &userPb.CheckUserRequest{
			UserName: loginVar.UserName,
			PassWord: loginVar.PassWord,
		}
		Uuid, err := rpc.CheckUser(ctx, req)
		if err != nil {
			handlers.SendResponse(c, errno.AuthorizationFailedErr, nil)
			c.Abort()
			return
		}
		c.Set(constants.IdentityKey, Uuid)
	}
}
