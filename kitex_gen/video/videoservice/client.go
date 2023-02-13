// Code generated by Kitex v0.4.4. DO NOT EDIT.

package videoservice

import (
	"context"
	video "first/kitex_gen/video"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Upload(ctx context.Context, Req *video.PublishListRequest, callOptions ...callopt.Option) (r *video.PublishListResponse, err error)
	GetVideoList(ctx context.Context, Req *video.GetVideoListRequest, callOptions ...callopt.Option) (r *video.GetVideoListResponse, err error)
	LikeVideo(ctx context.Context, Req *video.LikeVideoRequest, callOptions ...callopt.Option) (r *video.LikeVideoResponse, err error)
	IncrComment(ctx context.Context, Req *video.IncrCommentRequest, callOptions ...callopt.Option) (r *video.IncrCommentResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kVideoServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kVideoServiceClient struct {
	*kClient
}

func (p *kVideoServiceClient) Upload(ctx context.Context, Req *video.PublishListRequest, callOptions ...callopt.Option) (r *video.PublishListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Upload(ctx, Req)
}

func (p *kVideoServiceClient) GetVideoList(ctx context.Context, Req *video.GetVideoListRequest, callOptions ...callopt.Option) (r *video.GetVideoListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideoList(ctx, Req)
}

func (p *kVideoServiceClient) LikeVideo(ctx context.Context, Req *video.LikeVideoRequest, callOptions ...callopt.Option) (r *video.LikeVideoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.LikeVideo(ctx, Req)
}

func (p *kVideoServiceClient) IncrComment(ctx context.Context, Req *video.IncrCommentRequest, callOptions ...callopt.Option) (r *video.IncrCommentResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.IncrComment(ctx, Req)
}
