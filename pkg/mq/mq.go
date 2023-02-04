package mq

import (
	"context"
)

// 负责定义接口
type Subscriber interface {
	Subscribe(ctx context.Context) (ch chan *carpb.CarEntity, cleanUp func(), err error)
}

type Publisher interface {
	Publish(context.Context, *carpb.CarEntity) error
}
