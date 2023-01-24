package handlers

import (
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type Token struct {
	Token string `json:"token" query:"token"`
}
type UserId struct {
	UserId int64 `json:"user_id" query:"user_id"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func BuildResponse(err error) Response {
	Err := errno.ConvertErr(err)
	return Response{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
	}
}

// SendResponse pack response
func SendResponse(c *app.RequestContext, err error) {
	c.JSON(consts.StatusOK, BuildResponse(err))
}

func PackUser(u *user.User) *User {
	return &User{
		Id:            u.Id,
		Name:          u.UserName,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      u.IsFollow,
	}

}
