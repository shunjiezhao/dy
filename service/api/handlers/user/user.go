package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

//TODO: 将user rpc 中用到的接口提取出来, 方便mock

// Register 注册用户
// @tokenGenerator: 生成 token
// @RpcRegister: 调用 rpc 返回 userid
func (s *Service) Register(tokenGenerator func(data interface{}) (string, time.Time, error),
	RpcRegister func(ctx context.Context, req *userPb.RegisterRequest) (int64, error)) func(context2 *gin.Context) {
	return func(c *gin.Context) {
		var (
			param  RegisterRequest
			err    error
			req    *userPb.RegisterRequest
			token  string
			userId int64
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

		req = &userPb.RegisterRequest{
			UserName: param.UserName,
			PassWord: param.PassWord,
		}
		if RpcRegister == nil {
			RpcRegister = s.rpc.Register
		}
		userId, err = RpcRegister(c, req) // 方便mock
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}
		token, _, err = tokenGenerator(userId)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}
		SendRegisterResponse(c, userId, token)
		return

	errHandler:
		c.Abort()

	}
}
func (s *Service) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			loginVar LoginRequest
			err      error
			req      *userPb.CheckUserRequest
			Uuid     int64
		)
		notValid := func() bool {
			return len(loginVar.UserName) == 0 || len(loginVar.PassWord) == 0
		}
		err = c.ShouldBindQuery(&loginVar)
		if err != nil {
			err = c.ShouldBind(&loginVar)
		}

		if err != nil || notValid() {
			loginVar.UserName = c.Query("username")
			loginVar.PassWord = c.Query("password")
			if notValid() {
				handlers.SendResponse(c, errno.ParamErr)
				goto errHandler
			}

		}

		req = &userPb.CheckUserRequest{
			UserName: loginVar.UserName,
			PassWord: loginVar.PassWord,
		}

		Uuid, err = s.rpc.CheckUser(c, req)
		if Uuid == -1 {
			handlers.SendResponse(c, errno.AuthorizationFailedErr)
			goto errHandler

		}
		if err != nil {
			handlers.SendResponse(c, errno.ServiceErr)
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
			param    GetInfoRequest
			userID   string
			err      error
			req      *userPb.GetUserRequest
			userInfo *userPb.User
		)
		userID = c.Query("user_id")
		err = c.ShouldBindQuery(&param)
		if err != nil || len(param.GetToken()) == 0 || param.GetUserId() == 0 {
			err = c.ShouldBind(&param)
		}
		if (err != nil || param.GetUserId() <= int64(0)) && len(userID) == 0 {
			handlers.SendResponse(c, errno.ParamErr)
			goto errHandler
		}

		if len(userID) != 0 && param.GetUserId() == 0 {
			id, err := strconv.ParseInt(userID, 10, 64)
			if err != nil {
				handlers.SendResponse(c, errno.ParamErr)
				goto errHandler
			}
			param.SetUserId(id)
		}

		// 发送查询请求
		req = &userPb.GetUserRequest{Id: param.GetUserId()}
		userInfo, err = s.rpc.GetUserInfo(c, req)
		if err != nil {
			handlers.SendResponse(c, err)
			goto errHandler

		}
		SendGetInfoResponse(c, handlers.PackUser(userInfo))
		return

	errHandler:
		c.Abort()
	}
}
