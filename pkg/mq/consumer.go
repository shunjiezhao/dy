package mq

import (
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch        *amqp.Channel
	exchange  string
	queueName string
	key       string
}

//NewConsumer 	key 就是我们的 queue 名字AutoAck
func NewConsumer(dial *amqp.Connection, ex string, queueName string, key string, srvName string) *Consumer {
	if srvName == "" {
		srvName = "视频服务"
	}
	println(ex, queueName, key)
	channel, err := dial.Channel()
	if err != nil {
		panic(err)
	}

	p := &Consumer{
		ch:        channel,
		exchange:  ex,
		key:       key,
		queueName: queueName,
	}

	err = channel.ExchangeDeclare(
		p.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		klog.Errorf("无法创建交换机:%v", err)
		return nil
	}
	q, err := channel.QueueDeclare(
		srvName+": "+key, // name
		false,            // durable
		false,            // delete when unused
		true,             // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	p.queueName = q.Name
	if err != nil {
		klog.Errorf("无法创建Queue :%v", err)
	}
	err = p.ch.QueueBind(
		q.Name,     // queue name
		key,        // routing key
		p.exchange, // exchange
		false,
		nil)

	if err != nil {
		klog.Errorf("无法绑定Queue :%v", err)
		return nil
	}
	return p
}

func (p *Consumer) Consumer() (<-chan amqp.Delivery, error) {

	println("消费者: " + p.key)
	return p.ch.Consume(
		p.queueName,   // queue
		"消费者: "+p.key, // consumer
		true,          // auto ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // args
	)

}
