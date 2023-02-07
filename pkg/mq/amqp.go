package mq

import (
	"context"
	"first/pkg/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
	key      string
}

//NewPublisher key 就是我们的 queue 名字
func NewPublisher(dial *amqp.Connection, ex string, key string) *Publisher {

	channel, err := dial.Channel()
	if err != nil {
		panic(err)
	}

	p := &Publisher{
		ch:       channel,
		exchange: ex,
		key:      key,
	}
	p.NewTopicQueueProduct()
	return p
}
func (p *Publisher) NewTopicQueueProduct() {
	err := p.ch.ExchangeDeclare(
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
		return
	}

}

func (p *Publisher) Publish(ctx context.Context, data []byte) error {
	err := p.ch.PublishWithContext(ctx, p.exchange, p.key, false, false, amqp.Publishing{
		Body: data,
	})
	if err != nil {
		klog.Errorf("发送出错:%v", err)
		return errno.ServiceErr
	}
	return nil
}
