package user

import (
	"first/service/api/handlers"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

//TODO: 参数检验
type RegisterRequest struct {
	UserName string `json:"username" query:"username" vd:"len($) > 0"`
	PassWord string `json:"password" query:"password" vd:"len($) > 5'"`
}

type RegisterResponse struct {
	handlers.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type LoginResponse struct {
	handlers.Response
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func SendRegisterResponse(c *app.RequestContext, userId int64, token string) {
	c.JSON(consts.StatusOK, RegisterResponse{
		Response: handlers.Response{},
		UserId:   userId,
		Token:    token,
	})
}