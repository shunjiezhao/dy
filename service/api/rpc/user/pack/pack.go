package pack

import (
	"first/kitex_gen/user"
	"first/pkg/constants"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

func User(u *user.User) *handlers.User {
	if u == nil {
		return nil
	}
	return &handlers.User{
		Id:            u.Id,
		Name:          u.UserName,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      u.IsFollow,
	}
}
func Users(u []*user.User) []*handlers.User {
	users := make([]*handlers.User, 0)
	if len(u) == 0 {
		return users
	}

	for i := 0; i < len(u); i++ {
		users = append(users, User(u[i]))
	}
	return users
}

func Comments(com []*user.Comment) []*handlers.Comment {
	comments := make([]*handlers.Comment, 0)
	if len(com) == 0 {
		return comments
	}

	for i := 0; i < len(com); i++ {
		comments = append(comments, Comment(com[i]))
	}
	return comments
}

func Comment(c *user.Comment) *handlers.Comment {
	if c == nil {
		return nil
	}
	return &handlers.Comment{
		Id:         c.Id,
		User:       User(c.User),
		Content:    c.Content,
		CreateDate: time.Unix(c.CreateDate, 0).Format(constants.TimeFormatS),
	}
}
func FriendUsers(u []*user.FriendUser) []*handlers.FriendUser {
	users := make([]*handlers.FriendUser, len(u))
	if len(u) == 0 {
		return users
	}

	for i := 0; i < len(u); i++ {
		users[i] = FriendUser(u[i])
		klog.Infof("[pack.FriendList]: result: %v", users[i])

	}
	return users
}

func FriendUser(friendUser *user.FriendUser) *handlers.FriendUser {
	if friendUser == nil {
		return nil
	}
	return &handlers.FriendUser{
		User:    User(friendUser.User),
		Message: friendUser.Message,
		MsgType: friendUser.MsgType,
	}
}
