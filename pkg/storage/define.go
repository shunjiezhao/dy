package storage

import (
	"context"
)

type (
	AccessUrl string
	//Subscriber 消费者
	Subscriber interface {
		Subscribe(ctx context.Context) (ch chan AccessUrl, cleanUp func(), err error)
	}
	//Publisher 生产者
	Publisher interface {
		Publish(context.Context, AccessUrl) error
	}

	Info struct {
		Data  []byte `json:"data,omitempty"`
		Time  int64  `json:"time,omitempty"`
		Uuid  int64  `json:"uuid"`
		Title string `json:"title"`
		Hash  string `json:"hash"`
	}

	Storage interface {
		UploadFile(*Info) (string, string, error) // 返回我们的 获取链接
	}

	StorageFactory interface {
		Factory() Storage
	}
)
