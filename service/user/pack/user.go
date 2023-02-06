package pack

import (
	"errors"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/model/db"
	"github.com/cloudwego/kitex/pkg/klog"
)

// BuildBaseResp build baseResp from error
func BuildBaseResp(err error) *user.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}
func baseResp(err errno.ErrNo) *user.BaseResp {
	return &user.BaseResp{StatusCode: err.ErrCode, StatusMsg: err.ErrMsg}
}

func Users(dUsers []*db.User) []*user.User {
	users := make([]*user.User, len(dUsers))
	for i := 0; i < len(dUsers); i++ {
		users[i] = User(dUsers[i])
		klog.Infof("获取到信息, %#v", users[i])
	}
	return users
}
func User(dUser *db.User) *user.User {
	return &user.User{
		Id:            dUser.Uuid,
		UserName:      dUser.UserName,
		FollowCount:   dUser.FollowCount,
		FollowerCount: dUser.FollowerCount,
		IsFollow:      dUser.IsFollow,
	}
}
func Comment(com *db.Comment) *user.Comment {
	return &user.Comment{
		Id:         com.Id,
		User:       User(&com.User),
		Content:    com.Content,
		CreateDate: com.CreatedAt.Unix(), // 返回时间戳
	}
}

func Comments(com []*db.Comment) []*user.Comment {
	res := make([]*user.Comment, 0)
	for i := 0; i < len(com); i++ {
		res = append(res, Comment(com[i]))
	}
	return res
}
func Message(ms *db.Message) *user.Message {
	return &user.Message{
		FromUserId: ms.FromUserUuid,
		ToUserId:   ms.ToUserUuid,
		Content:    ms.Content,
		CreatedAtS: ms.CreatedAt.Unix(),
		MessageId:  ms.Id,
	}
}

func Messages(ms []*db.Message) []*user.Message {
	res := make([]*user.Message, len(ms))
	for i := 0; i < len(ms); i++ {
		res[i] = Message(ms[i])
	}
	return res
}
