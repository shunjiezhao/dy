package mq

import (
	"first/pkg/constants"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
)

func GetSaveVideoQueueName(id int64) string {
	return constants.SaveVideoPrefix + strconv.FormatInt(id%constants.VideoQCount, 10)
}
func GetSaveVideoQueueKey(id int64) string {
	return constants.SaveVideoKey + strconv.FormatInt(id%constants.VideoQCount, 10)
}

func GetSaveVideoIdx(id int64) int {
	return int(id % constants.VideoQCount)

}

func GetMqConnection() *amqp.Connection {
	var err error
	var conn *amqp.Connection

	conn, err = amqp.Dial(constants.MQConnURL)
	if err != nil {
		panic(err)
	}
	return conn
}

func UGetActionCommentQueueName(id int64) string {
	return constants.UActionCommentPrefix + strconv.FormatInt(id%constants.UActionCommentQCount, 10)
}
func UGetActionCommentQueueKey(id int64) string {
	return constants.UActionCommentKey + strconv.FormatInt(id%constants.UActionCommentQCount, 10)
}

type (
	ActionCommentInfo struct {
		Uuid        int64  `json:"uuid"`          // 用户id
		VideoId     int64  `json:"video_id"`      // 视频id
		ActionType  int32  ` json:"action_type"`  // 1-发布评论，2-删除评论
		CommentText string ` json:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
		CommentId   int64  ` json:"comment_id"`   // 要删除的评论id，在action_type=2的时候使用
	}
)

func UGetActionCommentIdx(id int64) int {
	return int(id % constants.UActionCommentQCount)
}

func VGetActionVideoComQueueName(id int64) string {
	return constants.UActionCommentPrefix + strconv.FormatInt(id%constants.VActionVideoComCountQCount, 10)
}
func VGetActionVideoComCountQueueKey(id int64) string {
	return constants.VActionVideoComCountKey + strconv.FormatInt(id%constants.VActionVideoComCountQCount, 10)
}

func VGetActionVideoComCountIdx(id int64) int {
	return int(id % constants.VActionVideoComCountQCount)
}

// NewSubConsumer 创建消费队列
func NewSubConsumer(count int64, ex string, getQueueName func(int64) string, getKey func(int64) string, srvName string) []*Consumer {
	// 创建publisher
	consumers := make([]*Consumer, count)
	connection := GetMqConnection()
	for i := 0; i < int(count); i++ {
		consumers[i] = NewConsumer(connection, ex, getQueueName(int64(i)), getKey(int64(i)), srvName)

	}
	return consumers
}
