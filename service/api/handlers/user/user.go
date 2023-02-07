package user

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/gin-gonic/gin"
)

// Register 注册用户
// @tokenGenerator: 生成 token
// @RpcRegister: 调用 rpc 返回 userid
func (s *Service) Register() func(context2 *gin.Context) {
	return func(c *gin.Context) {
		var (
			param  common.RegisterRequest
			err    error
			userId int64
			ctx    context.Context = c.Request.Context() // 方便 mock gin.Context 不是 context.Context
			data   middleware.PayLoad
		)
		err = c.ShouldBindQuery(&param)
		if err != nil {
			err = c.ShouldBind(&param)
		}
		// 参数校验
		if err != nil || len(param.UserName) == 0 || len(param.PassWord) == 0 {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		userId, err = s.rpc.Register(ctx, &param) // 方便mock
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}

		if userId <= 0 {
			handlers.SendResponse(c, errno.ServiceErr)
			goto errHandler

		}

		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}

		data = map[string]interface{}{
			constants.IdentityKey: userId,
			constants.UserNameKey: param.UserName,
		}
		c.Set(constants.PayLodKey, data)
		c.Set(constants.IdentityKey, userId)
		return

	errHandler:
		c.Abort()
	}
}
func (s *Service) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			param common.LoginRequest
			err   error
			user  *handlers.User
			ctx   context.Context = c.Request.Context()
			data  middleware.PayLoad
		)
		notValid := func() bool {
			return len(param.UserName) == 0 || len(param.PassWord) == 0
		}
		err = c.ShouldBind(&param)
		if err != nil {
			err = c.ShouldBindQuery(&param)
		}

		if err != nil || notValid() {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler

		}

		user, err = s.rpc.CheckUser(ctx, &param)
		if user == nil {
			handlers.SendResponse(c, errno.AuthorizationFailedErr)
			goto errHandler

		}
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}
		data = map[string]interface{}{
			constants.IdentityKey: user.Id,
			constants.UserNameKey: user.Name,
		}
		c.Set(constants.PayLodKey, data)
		c.Set(constants.IdentityKey, user.Id)
		return

	errHandler:
		c.Abort()
	}
}

func (s *Service) GetInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 检查参数
		var (
			param    common.GetInfoRequest
			err      error
			userInfo *handlers.User
			ctx      context.Context = c.Request.Context()
		)
		err = c.ShouldBindQuery(&param)
		if err != nil || len(param.GetToken()) == 0 || param.GetUserId() == 0 {
			err = c.ShouldBind(&param)
		}
		if err != nil {
			goto ParamErr
		}

		// 发送查询请求
		userInfo, err = s.rpc.GetUserInfo(ctx, param.UserId)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}

		common.SendGetInfoResponse(c, userInfo)
		return
	ParamErr:
		handlers.SendResponse(c, errno.ParamErr)
	errHandler:
		c.Abort()
	}
}
