package user

import (
	"context"
	"first/kitex_gen/user"
	"first/kitex_gen/user/chatservice"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	"first/pkg/util"
	"first/service/api/handlers"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	etcd "github.com/kitex-contrib/registry-etcd"
	"time"
)

var chatClient chatservice.Client

func InitChatRpc() {
	var err error
	resolver, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	clientName := "api-client"
	suite, _ := util.Trace(clientName)

	chatClient, err = chatservice.NewClient(
		constants.ChatServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithResolver(resolver), // etcd
		client.WithSuite(suite),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: clientName}),
	)
	if err != nil {
		panic(err)
	}
}

//go:generate mockgen -destination=../mock/user/chat.go -package=mock first/service/api/rpc/user ChatProxy
type ChatProxy interface {
	Save(context.Context, *handlers.Message) error
	GetList(context.Context, handlers.FromUserId, handlers.ToUserId) ([]*handlers.Message, error)
	GetFriendChatList(context.Context, handlers.FromUserId) ([]*handlers.Message, error)
}

type ChatRpcProxy struct {
	chatClient chatservice.Client
}

func (c *ChatRpcProxy) GetFriendChatList(ctx context.Context, id handlers.FromUserId) ([]*handlers.Message, error) {
	if id.UserId == 0 {
		return nil, errno.ParamErr

	}
	resp, err := c.chatClient.GetFriendChatList(ctx, &user.GetFriendChatRequest{
		FromUserId: id.GetFromUserId(),
	})
	if err != nil {
		klog.Errorf("[好友列表服务]: rpc 调用失败: %v", err)
		return nil, errno.RemoteErr
	}

	if resp == nil || respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return PackMsgS(resp.Msg), nil
}

func (c *ChatRpcProxy) Save(ctx context.Context, message *handlers.Message) error {
	resp, err := c.chatClient.SendMsg(ctx, &user.SaveMsgRequest{
		FromUserId: message.GetFromUserId(),
		ToUserId:   message.GetToUserId(),
		Content:    message.Content,
		CreatedAtS: time.Now().Unix(), //单位是秒
	})
	if err != nil {
		klog.Errorf("[消息服务]: rpc 调用失败: %v", err)
		return errno.RemoteErr
	}

	if resp == nil || respIsErr(resp.Resp) {
		return errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return nil
}

func (c *ChatRpcProxy) GetList(ctx context.Context, from handlers.FromUserId, to handlers.ToUserId) ([]*handlers.Message, error) {
	resp, err := c.chatClient.GetChatList(ctx, &user.GetChatListRequest{
		FromUserId: from.GetFromUserId(),
		ToUserId:   to.GetToUserId(),
	})

	if err != nil {
		klog.Errorf("[消息服务]: rpc 调用失败: %v", err)
		return nil, errno.RemoteErr
	}

	if respIsErr(resp.Resp) {
		return nil, errno.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMsg)
	}
	return PackMsgS(resp.MessageList), nil
}

func NewChatRpcProxy() ChatProxy {
	return &ChatRpcProxy{chatClient: chatClient}
}

func PackMsgS(MessageList []*user.Message) []*handlers.Message {
	res := make([]*handlers.Message, len(MessageList))
	for i := 0; i < len(MessageList); i++ {
		res[i] = &handlers.Message{
			Id:         MessageList[i].MessageId,
			ToUserId:   handlers.ToUserId{UserId: MessageList[i].ToUserId},
			FromUserId: handlers.FromUserId{UserId: MessageList[i].FromUserId},
			Content:    MessageList[i].Content,
			CreateTime: time.Unix(MessageList[i].CreatedAtS, 0).Format(constants.TimeFormatS),
		}
	}
	return res
}
