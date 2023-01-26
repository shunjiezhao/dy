package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetFollowUserList(t *testing.T) {
	// 删除
	assert.Empty(t, DB.Unscoped().Delete(&User{}, 1, 2, 3, 4).Error)

	// 数据
	//1. 创建 4 个 user
	users := insertUserHelper(4)
	for _, user := range users {
		DB.Unscoped().Delete(&user, user.Uuid)
		DB.Unscoped().Where("from_user_uuid = ? or to_user_uuid = ? ", user.Uuid, user.Uuid).Delete(&Follow{})
	}
	DB.CreateInBatches(users, 4) // 插入
	var input = `1->2 1->3 2->3 2->4 3->4 4->2 4->3`
	follows := getFollowsHelper(input)

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

func insertUserHelper(n int) []*User {
	var users []*User
	for i := 0; i < n; i++ {
		users = append(users, &User{
			Uuid:     int64(i + 1),
			UserName: fmt.Sprintf("user-%d", i+1),
			Password: "123456",
			NickName: fmt.Sprintf("user-%d", i+1),
		})
	}
	return users
}

func TestFollowUser(t *testing.T) {
	//1.创建数据
	Init()
	//1. 创建 3 个 user
	users := insertUserHelper(3)
	for _, user := range users {
		DB.Unscoped().Delete(&user, user.Uuid)
		DB.Unscoped().Where("from_user_uuid = ? or to_user_uuid = ? ", user.Uuid, user.Uuid).Delete(&Follow{})
	}
	DB.CreateInBatches(users, 4) // 插入

	var cases = []struct {
		name      string
		from      int64
		to        int64
		shouldErr bool
		// from follow , follower info,  to follow follower info
		shouldFollowInfo map[int][2]int64
		op               func(ctx context.Context, id int64, followerId int64) (bool, error)
	}{
		{
			name: "关注存在的用户",
			from: 1,
			to:   2,
			shouldFollowInfo: map[int][2]int64{
				1: [2]int64{0: 1, 1: 0},
				2: [2]int64{0: 0, 1: 1},
			},
			op: FollowUser,
		},
		{
			name:      "关注不存在的用户",
			shouldErr: true,
			from:      1,
			to:        0,
			shouldFollowInfo: map[int][2]int64{
				1: [2]int64{0: 1, 1: 0},
			},
			op: FollowUser,
		},
		{
			name:      "取消关注 没关注的",
			from:      1,
			to:        3,
			shouldErr: true,
			shouldFollowInfo: map[int][2]int64{
				1: [2]int64{0: 1, 1: 0},
				3: [2]int64{0: 0, 1: 0},
			},
			op: UnFollowUser,
		},
		{
			name: "再次关注",
			from: 1,
			to:   2,
			shouldFollowInfo: map[int][2]int64{
				1: [2]int64{0: 1, 1: 0},
				2: [2]int64{0: 0, 1: 1},
			},
			shouldErr: true,
			op:        FollowUser,
		},
		{
			name: "取消关注",
			from: 1,
			to:   2,
			shouldFollowInfo: map[int][2]int64{
				1: [2]int64{0: 0, 1: 0},
				2: [2]int64{0: 0, 1: 0},
			},
			op: UnFollowUser,
		},
	}

	ctx := context.Background()
	assert := assert.New(t)
	assert.NotEmpty(DB, "database not init")
	for _, c := range cases {
		_, err := c.op(ctx, c.to, c.from)
		var fromUser User
		dbErr := DB.First(&fromUser, c.from).Error
		assert.Empty(dbErr, "%s: get error: %s", c.name, dbErr)
		if c.shouldErr {
			assert.NotEmpty(err, "%s: should get error but not;", c.name)
			//检查 user 1 的信息
			ints := c.shouldFollowInfo[int(fromUser.Uuid)]
			assert.Equal(ints[0], fromUser.FollowCount, "%s: want follow count: %d; but %d;", c.name, ints[0], fromUser.FollowCount)
			assert.Equal(ints[1], fromUser.FollowerCount, "%s: want follower count: %d; but %d;", c.name, ints[1],
				fromUser.FollowerCount)
			continue
		}
		assert.Empty(err, "%s: get error: %s", c.name, err)
		// 查询信息
		var toUser User
		dbErr = DB.First(&toUser, c.to).Error
		assert.Empty(dbErr, "%s: get error: %s", c.name, dbErr)
		ints := c.shouldFollowInfo[int(toUser.Uuid)]
		assert.Equal(ints[0], toUser.FollowCount, "%s: want follow count: %d; but %d;", c.name, ints[0], toUser.FollowCount)
		assert.Equal(ints[1], toUser.FollowerCount, "%s: want follower count: %d; but %d;", c.name, ints[1],
			toUser.FollowerCount)
	}
	DB.Where("password = 123456").Unscoped().Delete(&User{})
	for _, user := range users {
		DB.Unscoped().Where("from_user_uuid = ? or to_user_uuid = ? ", user.Uuid, user.Uuid).Delete(&Follow{})
	}
}

func TestMain(t *testing.M) {
	//TODO: DOCKER run
	os.Exit(t.Run())
}
func getFollowsHelper(input string) []*Follow {
	split := strings.Split(input, " ")
	follows := make([]*Follow, len(split))
	for i, s := range split {
		ids := strings.Split(s, "->")
		from, _ := strconv.ParseInt(ids[0], 10, 64)
		to, _ := strconv.ParseInt(ids[1], 10, 64)
		follows[i] = &Follow{
			FromUserUuid: from,
			ToUserUuid:   to,
		}
	}
	return follows
}
