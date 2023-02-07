package storage

import (
	"bytes"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io/fs"
	"io/ioutil"
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
func (svc defaultOssStorage) UploadFile(info *Info) (string, string) {

	key := strconv.Itoa(int(info.Uuid)) + time.Now().Format(time.Kitchen) + ".mp4" // 拼接 filename
	ctx, cancelFunc := context.WithTimeout(context.Background(), svc.maxUploadTime*time.Second)
	defer cancelFunc()
	reader := bytes.NewReader(info.Data)
	ioutil.WriteFile("./save.mp4", info.Data, fs.ModePerm)
	_, err := svc.client.Object.Put(ctx, key, reader, nil)

	if err != nil {
		klog.Errorf("上传文件失败", err.Error())
		//TODO:继续放入消息队列中
		return "", ""
	}

	playUrl := svc.url + "/" + key
	if err != nil {
		klog.Errorf("上传文件失败", err.Error())
		return "", ""
	}

	return playUrl, playUrl + "?ci-process=snapshot&time=0&format=jpg"
}
