package service

import (
	"first/pkg/constants"
	"first/pkg/mq"
)

func NewSubConsumer() []*mq.Consumer {

	// 创建publisher
	consumers := make([]*mq.Consumer, constants.VideoQCount)
	connection := mq.GetMqConnection()
	for i := 0; i < int(constants.VideoQCount); i++ {
		consumers[i] = mq.NewConsumer(connection, constants.SaveVideoExName,
			mq.GetSaveVideoQueueName(int64(i)), mq.GetSaveVideoQueueKey(int64(i)))

	}
	return consumers
}
