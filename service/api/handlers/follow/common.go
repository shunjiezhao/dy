package follow

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

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
