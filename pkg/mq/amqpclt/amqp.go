package amqpclt

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

////Publish publishes a message.
//func (p *Publisher) Publish(c context.Context, car *carpb.CarEntity) error {
//	b, err := json.Marshal(car)
//	if err != nil {
//		return fmt.Errorf("cannot marshal: %v", err)
//	}
//
//	return p.ch.Publish(
//		p.exchange,
//		"",    //Key
//		false, //mandatory
//		false, //immedaiiote
//		amqp.Publishing{
//			Body: b,
//		},
//	)
//}
//func (s *Subscriber) Subscribe(ctx context.Context) (chan *carpb.CarEntity, func(), error) {
//	raw, cleanUp, err := s.SubscribeRaw(ctx)
//	if err != nil {
//		return nil, cleanUp, err
//	}
//	carCh := make(chan *carpb.CarEntity)
//	go func() {
//		for msg := range raw {
//			var car carpb.CarEntity
//			if err := json.Unmarshal(msg.Body, &car); err != nil {
//				s.logger.Error("can not unmarshal", zap.Error(err))
//			}
//			carCh <- &car
//		}
//		close(carCh)
//	}()
//	return carCh, cleanUp, nil
//}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	err = declareExchange(ch, exchange)

	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

type Subscriber struct {
	conn     *amqp.Connection
	exchange string
	logger   *zap.Logger
}

func NewSubscriber(conn *amqp.Connection, exchange string, logger *zap.Logger) (*Subscriber, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can not allocate channel:%v", err)
	}
	defer channel.Close()
	err = declareExchange(channel, exchange)
	if err != nil {
		return nil, fmt.Errorf("can not declare exchange:%v", err)
	}
	return &Subscriber{
		conn:     conn,
		exchange: exchange,
		logger:   logger,
	}, nil

}

func (s *Subscriber) SubscribeRaw(ctx context.Context) (<-chan amqp.Delivery, func(), error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, func() {}, fmt.Errorf("can not allocate channel:%v", err)
	}

	closeCh := func() {
		err := ch.Close()
		if err != nil {
			s.logger.Error("can not close channel", zap.Error(err))
		}
	}
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil)
	if err != nil {
		return nil, closeCh, fmt.Errorf("can not allocate queue:%v", err)
	}

	cleanUp := func() {
		// 最后关闭channel
		defer closeCh()
		// 先关闭 删除队列
		if _, err := ch.QueueDelete(q.Name, false, false, false); err != nil {
			s.logger.Error("can not delete queue", zap.String("queueName", q.Name))
		}
	}
	err = ch.QueueBind(
		q.Name,
		"",
		s.exchange,
		false,
		nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("can not  bind queue:%v", err)
	}
	consume, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("can not allocate consume:%v", err)
	}
	return consume, cleanUp, nil
}
