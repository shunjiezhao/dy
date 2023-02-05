package storage

import (
	"context"
	"mime/multipart"
	"time"
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

	Storage interface {
		UploadFile(string, *multipart.FileHeader, int64, time.Time) // 返回我们的 获取链接
	}

	StorageFactory interface {
		Factory() Storage
	}
)
