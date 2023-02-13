package storage

import (
	"bytes"
	"context"
	"first/pkg/errno"
	"first/pkg/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/tencentyun/cos-go-sdk-v5"
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
func (svc defaultOssStorage) UploadFile(info *Info) (string, string, error) {

	key := strconv.Itoa(int(info.Uuid)) + time.Now().Format(time.Kitchen) + ".mp4" // 拼接 filename
	ctx, cancelFunc := context.WithTimeout(context.Background(), svc.maxUploadTime*time.Second)
	defer cancelFunc()
	reader := bytes.NewReader(info.Data)
	// hash 检查
	hash := util.EncodeMD5(reader)
	if hash != info.Hash {
		klog.Errorf("hash 不想等;want: %s; but %s;", info.Hash, hash)
		return "", "", errno.VideoBrokeErr
	}

	_, err := svc.client.Object.Put(ctx, key, reader, nil)

	if err != nil {
		klog.Errorf("上传文件失败", err.Error())
		return "", "", errno.RemoteOssErr
	}

	playUrl := svc.url + "/" + key

	return playUrl, playUrl + "?ci-process=snapshot&time=0&format=jpg", nil
}
