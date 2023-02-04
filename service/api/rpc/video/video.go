package video

import (
	"context"
	videoPb "first/kitex_gen/video"
	"first/kitex_gen/video/videoservice"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var videoClient videoservice.Client

//go:generate mockgen -destination=../mock/male_mock.go -package=mock first/service/api/rpc/video RpcProxyIFace
type RpcProxyIFace interface {
	Upload(context.Context, *videoPb.PublishListRequest) error
	GetVideoList(ctx context.Context, Req *videoPb.GetVideoListRequest) ([]*videoPb.Video, error)
}

type RpcProxy struct {
	videoClient videoservice.Client
}

func (proxy RpcProxy) Upload(ctx context.Context, Req *videoPb.PublishListRequest) error {
	resp, err := proxy.videoClient.Upload(ctx, Req)
	if err != nil {
		return err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}

func (proxy RpcProxy) GetVideoList(ctx context.Context, Req *videoPb.GetVideoListRequest) ([]*videoPb.Video, error) {
	resp, err := proxy.videoClient.GetVideoList(ctx, Req)
	if err != nil {
		return nil, err
	}
	if resp.Resp != nil && resp.Resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return resp.VideoList, nil
}

func NewVideoProxy() RpcProxyIFace {
	return RpcProxy{videoClient: videoClient}
}

func InitApiVideoRpc() {
	var err error
	resolver, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	videoClient, err = videoservice.NewClient(
		constants.VideoServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithResolver(resolver), // etcd
	)

	if err != nil {
		panic(err)
	}
}
