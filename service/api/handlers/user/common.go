package user

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

//TODO: 参数检验
type RegisterRequest struct {
	UserName string `json:"username" query:"username" `
	PassWord string `json:"password" query:"password"`
}

type RegisterResponse struct {
	handlers.Response
	handlers.UserId
	handlers.Token
}

type LoginRequest struct {
	UserName string `json:"username"  query:"username"`
	PassWord string `json:"password" query:"password"`
}

type LoginResponse struct {
	handlers.Response
	Token string `json:"token"`
}

type GetInfoRequest struct {
	handlers.UserId
	handlers.Token
}
type GetInfoResponse struct {
	handlers.Response
	*handlers.User
}

func SendRegisterResponse(c *app.RequestContext, userId int64, token string) {
	c.JSON(consts.StatusOK, RegisterResponse{
		Response: handlers.BuildResponse(errno.Success),
		UserId:   handlers.UserId{UserId: userId},
		Token:    handlers.Token{Token: token},
	})
}
func SendGetInfoResponse(c *app.RequestContext, user *handlers.User) {
	c.JSON(consts.StatusOK, GetInfoResponse{
		Response: handlers.BuildResponse(errno.Success),
		User:     user,
	})
}
