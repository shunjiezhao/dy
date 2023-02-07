package video

import (
	"context"
	videoPb "first/kitex_gen/video"
	"first/kitex_gen/video/videoservice"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	"first/service/api/rpc/video/pack"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var videoClient videoservice.Client

//go:generate mockgen -destination=../mock/male_mock.go -package=mock first/service/api/rpc/video RpcProxyIFace
type RpcProxyIFace interface {
	Upload(ctx context.Context, param *videoPb.PublishListRequest) error
	GetVideoList(ctx context.Context, param *common.FeedRequest) ([]*handlers.Video, error)

	LikeVideo(ctx context.Context, param *videoPb.LikeVideoRequest) (err error)
}

type RpcProxy struct {
	videoClient videoservice.Client
}

func respIsErr(Resp *videoPb.VideoBaseResp) bool {
	return Resp != nil && Resp.StatusCode != errno.SuccessCode
}

func (proxy RpcProxy) LikeVideo(ctx context.Context, param *videoPb.LikeVideoRequest) (err error) {
	if param == nil {
		klog.Infof("请求参数为nil")
		return errno.ParamErr
	}

	resp, err := proxy.videoClient.LikeVideo(ctx, param)
	if err != nil {
		return errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}

func (proxy RpcProxy) Upload(ctx context.Context, param *videoPb.PublishListRequest) error {
	if param == nil {
		klog.Infof("请求参数为nil")
		return errno.ParamErr
	}

	resp, err := proxy.videoClient.Upload(ctx, param)
	if err != nil {
		return errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return nil
}

func (proxy RpcProxy) GetVideoList(ctx context.Context, param *common.FeedRequest) ([]*handlers.Video, error) {
	if param == nil {
		klog.Infof("请求参数为nil")
		return nil, errno.ParamErr
	}
	var req videoPb.GetVideoListRequest

	if param.GetAuthor {
		req.GetAuthor_ = true
		req.Author = param.Author
	} else if param.Uuid == 0 {
		req.TimeStamp = param.LatestTime // 获取当前点之后的
	} else {
		req.Uuid = param.Uuid
		req.IsLike = param.IsLike
	}

	resp, err := proxy.videoClient.GetVideoList(ctx, &req)
	if err != nil {
		return nil, errno.RemoteErr
	}
	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return pack.Videos(resp.VideoList), nil
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
