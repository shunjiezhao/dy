package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func TestComment(t *testing.T) {
	//1.创建数据
	if DB == nil {
		t.Fatalf("DB 初始化失败")
	}

	ctx := context.Background()
	//1. 创建 3 个 user
	users := insertUserHelper(3)
	var timeNow = time.Now()
	for _, user := range users {
		DB.Unscoped().Delete(&user, user.Uuid)
		DB.Unscoped().Where("from_user_uuid = ? or to_user_uuid = ? ", user.Uuid, user.Uuid).Delete(&Follow{})
	}

	DB.CreateInBatches(users, len(users)) // 插入
	commets := make([]*Comment, 3)
	// 插入 3 条评论
	for i := 0; i < 3; i++ {
		commets[i] = &Comment{
			Id:      int64(i + 1),
			Uuid:    users[i].Uuid,
			VideoId: 2,
			Content: users[i].UserName,
			Base: Base{
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
		}
		DB.Unscoped().Delete(commets[i], commets[i].Id)
		_, err := CreateComment(ctx, commets[i])
		if err != nil {
			t.Fatalf("插入评论失败")
		}

	}

	assert := assert.New(t)

	tests := []struct {
		name     string
		op       func() ([]*Comment, error)
		wantUser []*User
	}{
		{
			name: "获取用户1的评论",
			op: func() ([]*Comment, error) {
				return GetComment(ctx, 2, 1)
			},
			wantUser: []*User{users[0]},
		},
		{
			name: "获取用户2的评论",
			op: func() ([]*Comment, error) {
				return GetComment(ctx, 2, 1)
			},
			wantUser: []*User{users[1]},
		},
		{
			name: "获取所有的评论",
			op: func() ([]*Comment, error) {
				return GetCommentAfterTime(ctx, 2, time.Now().Add(time.Second).Unix(), 10)
			},
			wantUser: users,
		},
	}
	for _, test := range tests {
		comment, err := test.op()
		assert.Nilf(err, "获取评论得到错误: %v", err)
		assert.Equalf(len(test.wantUser), len(comment), "%s:获取评论数不一致", test.name)

		sort.Slice(comment, func(i, j int) bool {
			return comment[i].Uuid < comment[j].Uuid
		})
		for i := 0; i < len(comment); i++ {
			assert.Equalf(comment[i].User.UserName, users[i].UserName, "%s: 获取用户名不一致", test.name)
			assert.Equalf(comment[i].User.Uuid, users[i].Uuid, "%s: 获取用户Uuid不一致", test.name)
			assert.Equalf(comment[i].User.Password, users[i].Password, "%s: 获取用户密码不一致", test.name)
		}

	}
}
