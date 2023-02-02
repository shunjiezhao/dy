package user

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/rpc/user"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

//Service 用户微服务代理
type Service struct {
	rpc user.RpcProxyIFace
}

func New() *Service {
	return &Service{user.NewUserProxy()}
}

// User
//TODO: 参数检验
type (
	RegisterRequest struct {
		UserName string `json:"username" form:"username" `
		PassWord string `json:"password" form:"password"`
	}
	RegisterResponse struct {
		handlers.Response
		handlers.UserId
		handlers.Token
	}
	LoginRequest struct {
		UserName string `json:"username"  form:"username"`
		PassWord string `json:"password" form:"password"`
	}

	LoginResponse struct {
		handlers.Response
		Token string `json:"token"`
	}
	GetInfoRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetInfoResponse struct {
		handlers.Response
		*handlers.User `json:"user"`
	}
)

func SendRegisterResponse(c *gin.Context, userId int64, token string) {
	c.JSON(consts.StatusOK, RegisterResponse{
		Response: handlers.BuildResponse(errno.Success),
		UserId:   handlers.UserId{UserId: userId},
		Token:    handlers.Token{Token: token},
	})
}
func SendGetInfoResponse(c *gin.Context, user *handlers.User) {
	c.JSON(consts.StatusOK, GetInfoResponse{
		Response: handlers.BuildResponse(errno.Success),
		User:     user,
	})
}

// Follow

type (
	ActionType    int32
	ActionRequest struct {
		handlers.Token
		UserId     int64 `json:"to_user_id" form:"to_user_id"`
		ActionType `json:"action_type" form:"action_type"`
	}

	ActionResponse struct {
		handlers.Response
	}

	GetUserFollowerListRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetUserFollowerListResponse struct {
		handlers.Response
		Users []*handlers.User `json:"user_list,omitempty"`
	}

	GetUserFollowListRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetUserFollowListResponse struct {
		handlers.Response
		Users []*handlers.User `json:"users,omitempty"`
	}
)

func (a ActionType) IsFollow() bool {
	return a == 1
}
func SendUserListResponse(c *gin.Context, users []*handlers.User) {
	c.JSON(consts.StatusOK, GetUserFollowerListResponse{
		Response: handlers.BuildResponse(errno.Success),
		Users:    users,
	})
}
