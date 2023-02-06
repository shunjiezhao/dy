package pack

import (
	"first/kitex_gen/user"
	"first/service/api/handlers"
)

func User(u *user.User) *handlers.User {
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

func PackComments(com []*user.Comment) []*handlers.Comment {
	comments := make([]*handlers.Comment, 0)
	if len(com) == 0 {
		return comments
	}

	for i := 0; i < len(com); i++ {
		comments = append(comments, PackComment(com[i]))
	}
	return comments
}

func PackComment(c *user.Comment) *handlers.Comment {
	return &handlers.Comment{
		Id:         c.Id,
		User:       User(c.User),
		Content:    c.Content,
		CreateDate: c.CreateDate,
	}
}
