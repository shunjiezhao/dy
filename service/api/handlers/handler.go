package handlers

import (
	"first/pkg/constants"
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
	UserId struct {
		UserId int64 `json:"user_id" form:"user_id"`
	}
	ToUserId struct {
		UserId int64 `json:"to_user_id" form:"to_user_id"`
	}
	FromUserId struct {
		UserId int64 `json:"from_user_id" form:"from_user_id"`
	}
)

func (t Token) GetToken() string {
	return t.Token
}
func (t ToUserId) GetToUserId() int64 {
	return t.UserId
}
func (t FromUserId) GetFromUserId() int64 {
	return t.UserId
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
		Id            int64  `json:"id"`
		Author        *User  `json:"author"`
		PlayUrl       string `json:"play_url" json:"play_url"`
		CoverUrl      string `json:"cover_url"`
		FavoriteCount int64  `json:"favorite_count"`
		CommentCount  int64  `json:"comment_count"`
		IsFavorite    bool   `json:"is_favorite"`
	}
	Comment struct {
		Id         int64  `json:"id"`
		User       *User  `json:"user"`
		Content    string `json:"content"`
		CreateDate int64  `json:"create_date"`
	}

	Message struct {
		Id int64 `json:"id"` // 消息id
		ToUserId
		FromUserId
		Content    string `json:"content"`     // 消息内容
		CreateTime string `json:"create_time"` // 消息创建时间
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
func GetTokenUserId(c *gin.Context) int64 {
	claim := c.MustGet(constants.IdentityKey)

	var curUserId int64
	tmp, ok := claim.(float64)
	if ok {
		curUserId = int64(tmp)
	} else {
		curUserId = claim.(int64)
	}
	return curUserId
}
