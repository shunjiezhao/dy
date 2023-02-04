package handlers

import (
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

type (
	Response struct {
		StatusCode int64  `json:"status_code"`
		StatusMsg  string `json:"status_msg,omitempty"`
	}
	Token struct {
		Token string `json:"token" form:"token"`
	}
)

func (t Token) GetToken() string {
	return t.Token
}

type UserId struct {
	UserId int64 `json:"user_id" form:"user_id"`
}

func (t UserId) GetUserId() int64 {
	return t.UserId
}
func (t UserId) SetUserId(userId int64) {
	t.UserId = userId
}

type (
	User struct {
		Id            int64  `json:"id"`
		Name          string `json:"name"`
		FollowCount   int64  `json:"follow_count"`
		FollowerCount int64  `json:"follower_count"`
		IsFollow      bool   `json:"is_follow"`
	}
	Video struct {
		Id            int64  `json:"id,omitempty"`
		Author        User   `json:"author"`
		PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
		CoverUrl      string `json:"cover_url,omitempty"`
		FavoriteCount int64  `json:"favorite_count,omitempty"`
		CommentCount  int64  `json:"comment_count,omitempty"`
		IsFavorite    bool   `json:"is_favorite,omitempty"`
	}
)

func BuildResponse(err error) Response {
	Err := errno.ConvertErr(err)
	return Response{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
	}
}

// SendResponse pack response
func SendResponse(c *gin.Context, err error) {
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
func PackUsers(u []*user.User) []*User {
	users := make([]*User, 0)
	if len(u) == 0 {
		return users
	}

	for i := 0; i < len(u); i++ {
		users = append(users, PackUser(u[i]))
	}
	return users
}
