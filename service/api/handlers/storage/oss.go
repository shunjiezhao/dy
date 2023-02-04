package storage

import (
	"context"
	videoPb "first/kitex_gen/video"
	"first/pkg/logger"
	video2 "first/service/api/rpc/video"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DefaultOssFactory struct {
	Key string
	Id  string
	Url string
}

func (f DefaultOssFactory) Factory() Storage {
	u, _ := url.Parse(f.Url)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  f.Id,
			SecretKey: f.Key,
		},
	})
	return defaultOssStorage{
		client:        client,
		url:           f.Url,
		maxUploadTime: time.Millisecond * 5,
	}
}

type defaultOssStorage struct {
	client        *cos.Client
	url           string
	maxUploadTime time.Duration
}

//UploadFile 默认 5s 上传文件失败
func (svc defaultOssStorage) UploadFile(title string, fileHeader *multipart.FileHeader, uuid int64, t time.Time) {
	open, err := fileHeader.Open()

	if err != nil {
		return
	}

	key := strconv.Itoa(int(uuid)) + "-" + title + fileHeader.Filename // 拼接 filename
	ctx, cancelFunc := context.WithTimeout(context.Background(), svc.maxUploadTime*time.Second)
	defer cancelFunc()
	_, err = svc.client.Object.Put(ctx, key, open, nil)

	if err != nil {
		logger.GetLogger().Error("上传文件失败", zap.String("err", err.Error()))
		return
	}
	err = video2.NewVideoProxy().Upload(ctx, &videoPb.PublishListRequest{
		Author:  uuid,
		PlayUrl: svc.url + "/" + key,
		//TODO: 生成视频截图
		CoverUrl: "TODO",
		Title:    title,
	})
	if err != nil {
		logger.GetLogger().Error("上传文件失败", zap.String("error", err.Error()))
	}

	return
}
