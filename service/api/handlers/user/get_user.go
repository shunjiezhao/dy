package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/rpc"
	"github.com/cloudwego/hertz/pkg/app"
)

func Register(ctx context.Context, c *app.RequestContext) {
	var param RegisterRequest
	err := c.Bind(&param)
	if err != nil {
		handlers.SendResponse(c, err, nil)
		return
	}
	if len(param.UserName) == 0 || len(param.PassWord) == 0 {
		handlers.SendResponse(c, errno.ParamErr, nil)
		return
	}
	req := userPb.RegisterRequest{
		UserName: param.UserName,
		PassWord: param.PassWord,
	}
	userId, err := rpc.Register(ctx, &req)
	if err != nil {
		handlers.SendResponse(c, err, nil)
		return
	}
	//TODO: jwt 生成token
	SendRegisterResponse(c, userId, "this is test token")
}
func Login(ctx context.Context, c *app.RequestContext) {
	panic("not impl")
}
