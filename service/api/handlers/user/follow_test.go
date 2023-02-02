package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFollowUser(t *testing.T) {
	// 发送 Follow 请求
	// 1. 关注 自己
	//	2. 关注未关注的人
	//	3. 关注已关注的人
	//	4. 取消关注未关注的人
	//	5. 取消关注已关注的人
	e, _, _ := newExpect(t, serverAddr)
	userIda, tokenA := getTestUserToken(testUserA, e)
	if userIda == 0 || tokenA == "" {
		t.Fatal("not get token and userid")
	}
	userIdb, tokenB := getTestUserToken(testUserB, e)
	if userIdb == 0 || tokenB == "" {
		t.Fatal("not get token and userid")
	}

	route := "/douyin/relation/action"
	// 关注成功
	resp := e.POST(route).WithQuery("token", tokenA).WithQuery("to_user_id", userIdb).WithQuery("action_type", 1).
		WithFormField("token", tokenA).WithFormField("to_user_id", userIdb).WithFormField("action_type", 1).
		Expect().
		Status(http.StatusOK).JSON().
		Object()
	resp.Value("status_code").Number().Equal(0)
	route = "/douyin/relation/follow/list/"
	followListResp := e.GET(route).
		WithQuery("token", tokenA).WithQuery("user_id", userIda).
		WithFormField("token", tokenA).WithFormField("user_id", userIda).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	resp.Value("status_code").Number().Equal(0)

	// 包含 用户 B
	containTestUserB := false
	for _, element := range followListResp.Value("user_list").Array().Iter() {
		user := element.Object()
		user.ContainsKey("id")
		if int(user.Value("id").Number().Raw()) == userIdb {
			containTestUserB = true
		}
	}
	assert.True(t, containTestUserB, "Follow test user failed")
	route = "/douyin/relation/follower/list/"
	followerListResp := e.GET(route).
		WithQuery("token", tokenB).WithQuery("user_id", userIdb).
		WithFormField("token", tokenB).WithFormField("user_id", userIdb).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	resp.Value("status_code").Number().Equal(0)

	// 包含 用户 B
	containTestUserA := false
	for _, element := range followerListResp.Value("user_list").Array().Iter() {
		user := element.Object()
		user.ContainsKey("id")
		if int(user.Value("id").Number().Raw()) == userIda {
			containTestUserA = true
		}
	}
	assert.True(t, containTestUserA, "Follower test user failed")
}
