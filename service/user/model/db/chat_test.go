package db

import (
	"context"
	"first/pkg/constants"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	//1.创建数据
	if DB == nil {
		t.Fatalf("DB 初始化失败")
	}
	DB.Table(constants.MessageTableName).Where("from_user_uuid in ?", []int64{1, 2, 3}).Delete(&Message{})
	// 时间 依次减小
	// msg :  0     1    2 	   3    4     5    6
	input := "1->2 1->2 2->3 3->1 1->3 2->3 2->2"
	msgs := getChatHelper(input, 1000)
	//  "获取1的好友消息",
	DB.CreateInBatches(msgs, len(msgs))

	var all []*Message
	DB.Find(&all)

	list, err := GetFriendChatList(DB, context.Background(), 1)
	if err != nil {
		t.Fatalf("获取出错")
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].OtherId < list[j].OtherId
	})
	assert.Equal(t, 2, len(list), "长度不一样")
	assert.Equal(t, list[0].MySend, true)  // 2
	assert.Equal(t, list[1].MySend, false) // 3
	assert.Equal(t, list[0].OtherId, int64(2))
	assert.Equal(t, list[1].OtherId, int64(3))
	assert.Equal(t, list[0].Content, "0")
	assert.Equal(t, list[1].Content, "3")
}

func getChatHelper(input string, interval int64) []*Message {
	split := strings.Split(input, " ")
	unix := time.Now().Unix()
	msgs := make([]*Message, len(split))
	for i, s := range split {
		ids := strings.Split(s, "->")
		from, _ := strconv.ParseInt(ids[0], 10, 64)
		to, _ := strconv.ParseInt(ids[1], 10, 64)
		msgs[i] = &Message{
			Id:           int64(i + 1),
			FromUserUuid: from,
			ToUserUuid:   to,
			Content:      strconv.FormatInt(int64(i), 10),
			Base: Base{
				CreatedAt: time.Unix(unix-interval*int64(i), 0),
				UpdatedAt: time.Unix(unix-interval*int64(i), 0),
			},
		}
	}
	return msgs
}
