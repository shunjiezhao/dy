package mq

import (
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// RabbitConsumer rabbitmq 消费者
type RabbitConsumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	done          chan struct{}
	addr          string
	exchange      string
	queue         string
	routingKey    string
	consumerTag   string
	autoDelete    bool
	handler       func([]byte) error
	delivery      <-chan amqp.Delivery
}

// NewDelayConsumer 创建消费者
func NewDelayConsumer(config Config, handler func([]byte) error, Tag string) *RabbitConsumer {
	return &RabbitConsumer{
		addr:        config.Addr,
		exchange:    config.Exchange,
		queue:       config.Queue,
		routingKey:  config.RoutingKey,
		consumerTag: Tag,
		autoDelete:  config.AutoDelete,
		handler:     handler,
		done:        make(chan struct{}),
	}
}

func (c *RabbitConsumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}
	go c.ReConnect()
	return nil
}

func (c *RabbitConsumer) Stop() {
	close(c.done)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.consumerTag, true); err != nil {
			klog.Error("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			klog.Error("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

func (c *RabbitConsumer) Run() (err error) {
	if c.conn, err = amqp.Dial(c.addr); err != nil {
		return err
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		c.conn.Close()
		return err
	}

	defer func() {
		if err != nil {
			c.channel.Close()
			c.conn.Close()
		}
	}()

	// 声明一个主要使用的 exchange
	err = c.channel.ExchangeDeclare(
		c.exchange, "x-delayed-message", true, c.autoDelete, false, false, amqp.Table{
			"x-delayed-type": "fanout",
		})
	if err != nil {
		return err
	}

	// 声明一个延时队列, 延时消息就是要发送到这里
	q, err := c.channel.QueueDeclare(c.queue, false, c.autoDelete, false, false, nil)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(q.Name, "", c.exchange, false, nil)
	if err != nil {
		return err
	}

	c.delivery, err = c.channel.Consume(
		q.Name, c.consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go c.Handle()

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))
	return
}

func (c *RabbitConsumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				klog.Error("rabbitmq consumer - connection NotifyClose: ", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				klog.Error("rabbitmq consumer - channel NotifyClose: ", err)
			}
		case <-c.done:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				klog.Error("rabbitmq consumer - channel cancel failed: ", err)
			}

			if err := c.conn.Close(); err != nil {
				klog.Error("rabbitmq consumer - channel cancel failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			klog.Error(err)
		}
		for err := range c.connNotify {
			klog.Error(err)
		}

	quit:
		for {
			select {
			case <-c.done:
				return
			default:
				klog.Error("rabbitmq consumer - reconnect")

				if err := c.Run(); err != nil {
					klog.Error("rabbitmq consumer - failCheck: ", err)
					// sleep 15s reconnect
					time.Sleep(time.Second * 15)
					continue
				}
				break quit
			}
		}
	}
}

func (c *RabbitConsumer) Handle() {
	for d := range c.delivery {
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err != nil {
				// 重新入队，否则未确认的消息会持续占用内存，这里的操作取决于你的实现，你可以当出错之后并直接丢弃也是可以的
				_ = delivery.Reject(true)
			} else {
				_ = delivery.Ack(false)
			}
		}(d)
	}
}
