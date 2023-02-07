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
