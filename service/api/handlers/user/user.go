package user

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// Register 注册用户
// @tokenGenerator: 生成 token
// @RpcRegister: 调用 rpc 返回 userid
func (s *Service) Register(tokenGenerator func(data interface{}) (string, time.Time, error),
) func(context2 *gin.Context) {
	return func(c *gin.Context) {
		var (
			param  common.RegisterRequest
			err    error
			token  string
			userId int64
			ctx    context.Context = c.Request.Context() // 方便 mock gin.Context 不是 context.Context
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

		token, _, err = tokenGenerator(userId)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}
		common.SendRegisterResponse(c, userId, token)
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
			Uuid  int64
			ctx   context.Context = c.Request.Context()
		)
		notValid := func() bool {
			return len(param.UserName) == 0 || len(param.PassWord) == 0
		}
		err = c.ShouldBindQuery(&param)
		if err != nil {
			err = c.ShouldBind(&param)
		}

		if err != nil || notValid() {
			param.UserName = c.Query("username")
			param.PassWord = c.Query("password")
			if notValid() {
				handlers.SendResponse(c, errno.ParamErr)
				goto errHandler
			}

		}

		Uuid, err = s.rpc.CheckUser(ctx, &param)
		if Uuid == -1 {
			handlers.SendResponse(c, errno.AuthorizationFailedErr)
			goto errHandler

		}
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}

		c.Set(constants.IdentityKey, Uuid)
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
			userID   string
			err      error
			userInfo *handlers.User
			ctx      context.Context = c.Request.Context()
		)
		userID = c.Query("user_id")
		err = c.ShouldBindQuery(&param)
		if err != nil || len(param.GetToken()) == 0 || param.GetUserId() == 0 {
			err = c.ShouldBind(&param)
		}
		if (err != nil || param.GetUserId() <= int64(0)) && len(userID) == 0 {
			goto ParamErr
		}

		if len(userID) != 0 && param.GetUserId() == 0 {
			id, err := strconv.ParseInt(userID, 10, 64)
			if err != nil {
				goto ParamErr
			}
			param.SetUserId(id)
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
