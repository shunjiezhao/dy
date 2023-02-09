package handler

import (
	"context"
	userPb "first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/mq"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/json"
)

//UpdateVideoInfoConStart 异步上传视频的接受者
func (s *UserServiceImpl) UpdateVideoInfoConStart() func() {
	cons := mq.NewSubConsumer(constants.UActionCommentQCount, constants.UActionCommentExName, mq.UGetActionCommentQueueName, mq.UGetActionCommentQueueKey, "用户服务")
	done := make(chan struct{})
	fmt.Println(len(cons))
	for i := 0; i < len(cons); i++ {
		go func(i int, con *mq.Consumer) {
			consumer, err := con.Consumer()
			if err != nil {
				klog.Errorf("%d 号 消息队列挂掉, %v", i, err)
				return
			}
			klog.Infof("%d 号 消息队列启动", i)
			for data := range consumer {
				var info mq.ActionCommentInfo
				err = json.Unmarshal(data.Body, &info)
				if err != nil {
					klog.Infof("%d 好 unmarshal 失败 %v", i, err)
					return
				}

				klog.Infof("%d 号 消息队列获取到参数:%#v", i, info)

				for i := 0; i < 2; i++ {
					_, err = s.ActionComment(context.Background(), &userPb.ActionCommentRequest{
						Uuid:        info.Uuid,
						VideoId:     info.VideoId,
						ActionType:  info.ActionType,
						CommentText: info.CommentText,
						CommentId:   info.CommentId,
					})
					if err != nil {
						klog.Errorf("%d 号 消息队列处理失败: DB保存失败%v", i, err)
						continue
					}
					break
				}
				klog.Infof("%d 号 保存评论成功", i)

			}

		}(i, cons[i])
	}
	cleanUp := func() {
		done <- struct{}{}
	}
	return cleanUp

}
