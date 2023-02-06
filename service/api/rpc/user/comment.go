package user

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/api/handlers"
	"first/service/api/handlers/common"
	pack2 "first/service/api/rpc/user/pack"
	"github.com/cloudwego/kitex/pkg/klog"
)

func (proxy RpcProxy) ActionComment(ctx context.Context, param *common.CommentActionRequest) (r *handlers.Comment,
	err error) {
	req := userPb.ActionCommentRequest{
		Uuid:        param.UserId,
		VideoId:     param.VideoId,
		ActionType:  int32(param.CommentActionType),
		CommentText: param.CommentText,
		CommentId:   param.CommentId,
	}

	resp, err := proxy.userClient.ActionComment(ctx, &req)
	if err != nil || resp.Comment == nil {
		klog.Errorf("[UserRpc.ActionComment]: 失败")
		return nil, err
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return pack2.PackComment(resp.Comment), nil
}

func (proxy RpcProxy) GetComment(ctx context.Context, param *common.CommentListRequest) (r []*handlers.Comment, err error) {

	req := &userPb.GetCommentRequest{
		VideoId: param.VideoId,
	}

	resp, err := proxy.userClient.GetComment(ctx, req)
	if err != nil || resp == nil {
		klog.Errorf("[UserRpc.GetComment]: 失败")
		return nil, err
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}

	return pack2.PackComments(resp.Comment), nil
}
