package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
	"github.com/cloudwego/kitex/pkg/klog"
)

func (proxy RpcProxy) ActionComment(ctx context.Context, Req *userPb.ActionCommentRequest) (r *userPb.Comment, err error) {
	resp, err := proxy.userClient.ActionComment(ctx, Req)
	if err != nil || resp.Comment == nil {
		klog.Errorf("[UserRpc.ActionComment]: 失败")
		return nil, err
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return resp.Comment, nil
}

func (proxy RpcProxy) GetComment(ctx context.Context, Req *userPb.GetCommentRequest) (r []*userPb.Comment, err error) {
	resp, err := proxy.userClient.GetComment(ctx, Req)
	if err != nil || resp == nil {
		klog.Errorf("[UserRpc.GetComment]: 失败")
		return nil, err
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return resp.Comment, nil
}
