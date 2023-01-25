package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGetFollowUserList(t *testing.T) {
	Init()
	// 删除
	assert.Empty(t, DB.Unscoped().Delete(&User{}, 1, 2, 3, 4).Error)

	// 数据
	//1. 创建 4 个 user
	var users []*User
	for i := 0; i < 4; i++ {
		users = append(users, &User{
			Uuid:     int64(i + 1),
			UserName: fmt.Sprintf("user-%d", i+1),
			Password: "123456",
			NickName: fmt.Sprintf("user-%d", i+1),
		})
	}
	DB.CreateInBatches(users, 4) // 插入
	var follows []*Follow
	var i int64
	for ; i < 4; i++ {
		// 1->2, 2->3, 3->4
		follows = append(follows, &Follow{
			FromUserUuid: i + 1,
			ToUserUuid:   i + 2,
		})
	}
	for i = 0; i < 2; i++ {
		// 1->3, 2->4
		follows = append(follows, &Follow{
			FromUserUuid: i + 1,
			ToUserUuid:   i + 3,
		})
		// 4 -> 2 ,3
		follows = append(follows, &Follow{
			FromUserUuid: 4,
			ToUserUuid:   i + 2,
		})
	}

	DB.CreateInBatches(follows, len(follows))
	// 插入
	assert := assert.New(t)

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*5)
	cases := []struct {
		name         string
		shouldErr    bool
		shouldFollow [5]bool
		shouldLen    int
		op           func() ([]*User, error)
	}{
		{
			name:      "query user-1 follower list; should be empty",
			shouldErr: false,
			shouldLen: 0,
			op: func() ([]*User, error) {
				return GetFollowerUserList(ctx, 1)
			},
		},
		{
			name:      "query user-2 follower list; should get user-1,4",
			shouldErr: false,
			shouldLen: 2,
			shouldFollow: [5]bool{
				4: true,
			},
			op: func() ([]*User, error) {
				return GetFollowerUserList(ctx, 2)
			},
		},
		{
			name:         "query user-3 follower list; should get user-1,2,4; follow 4",
			shouldErr:    false,
			shouldLen:    3,
			shouldFollow: [5]bool{4: true},
			op: func() ([]*User, error) {
				return GetFollowerUserList(ctx, 3)
			},
		},
		{
			name:      "query user-4 follower list;should get user-2,3; and follow 2,3",
			shouldErr: false,
			shouldLen: 2,
			shouldFollow: [5]bool{
				2: true, 3: true,
			},
			op: func() ([]*User, error) {
				return GetFollowerUserList(ctx, 4)
			},
		},
	}

	for _, c := range cases {
		users, err := c.op()
		if c.shouldErr {
			assert.NotEmpty(err, "%s: should get error but not;", c.name)
		}
		if err != nil {
			assert.Empty(err, "%s: get error: %s", c.name, err.Error())
		}
		assert.Equal(c.shouldLen, len(users), "%s: ", c.name)
		for _, user := range users {
			assert.Equal(c.shouldFollow[user.Uuid], user.IsFollow, "%s: ", c.name)
		}
	}
}

func TestMain(t *testing.M) {
	os.Exit(t.Run())
}
