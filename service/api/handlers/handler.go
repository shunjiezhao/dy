package handlers

import (
	"first/kitex_gen/user"
	videoPb "first/kitex_gen/video"
	"first/pkg/constants"
	"first/pkg/errno"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/klog"
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
func PackVideos(videos []*videoPb.Video, users []*user.User, isOne bool) []*Video {
	var (
		one *User
		idx map[int64]int
	)
	if isOne {
		one = PackUser(users[0])
	} else {
		idx = make(map[int64]int, len(users))
		for i := 0; i < len(users); i++ { // 记录 用户 id 在 数组的位置
			idx[users[i].Id] = i
		}
	}

	res := make([]*Video, len(videos))
	for i := 0; i < len(videos); i++ {

		res[i] = &Video{
			Id:            videos[i].Id,
			PlayUrl:       videos[i].PlayUrl,
			CoverUrl:      videos[i].CoverUrl,
			FavoriteCount: videos[i].FavoriteCount,
			CommentCount:  videos[i].CommentCount,
			IsFavorite:    videos[i].IsFavorite,
		}

		if isOne {
			res[i].Author = one
			continue
		}

		if j, ok := idx[videos[i].Author]; ok {
			res[i].Author = PackUser(users[j])
		} else {
			klog.Infof("[Pack] 无法找到对应 id; users: %d Author: %d", j, videos[i].Author)
		}
	}
	return res

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
func PackComments(com []*user.Comment) []*Comment {
	comments := make([]*Comment, 0)
	if len(com) == 0 {
		return comments
	}

	for i := 0; i < len(com); i++ {
		comments = append(comments, PackComment(com[i]))
	}
	return comments
}

func PackComment(c *user.Comment) *Comment {
	return &Comment{
		Id:         c.Id,
		User:       PackUser(c.User),
		Content:    c.Content,
		CreateDate: c.CreateDate,
	}
}
