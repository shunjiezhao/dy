package main

import (
	"bytes"
	"context"
	"encoding/gob"
	userPb "first/kitex_gen/user"
	user "first/kitex_gen/user/userservice"
	video "first/kitex_gen/video"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/mq"
	"first/pkg/storage"
	"first/pkg/util"
	"first/service/video/pack"
	"first/service/video/service"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/json"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	userRpc user.Client
}

func (s *VideoServiceImpl) IncrComment(ctx context.Context, req *video.IncrCommentRequest) (resp *video.IncrCommentResponse, err error) {
	resp = new(video.IncrCommentResponse)
	if req == nil {
		goto ParamErr
	}
	err = service.NewVideoItemService(ctx).IncrCommentCount(req)
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.LikeOpErr)
		return resp, nil

	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
ParamErr:
	resp.Resp = pack.BuildBaseResp(errno.ParamErr)
	return
}

func (s *VideoServiceImpl) LikeVideo(ctx context.Context, req *video.LikeVideoRequest) (resp *video.LikeVideoResponse, err error) {
	resp = new(video.LikeVideoResponse)
	if req == nil {
		goto ParamErr
	}

	if req.VideoId == 0 || req.ActionType == nil {
		goto ParamErr

	}

	if req.ActionType.ActionType == 1 {
		err = service.NewLikeService(ctx).LikeVideo(req)
	} else {
		err = service.NewLikeService(ctx).UnLikeVideo(req) // 不喜欢
	}
	if err != nil {
		resp.Resp = pack.BuildBaseResp(errno.LikeOpErr)
		return resp, nil

	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return
ParamErr:
	resp.Resp = pack.BuildBaseResp(errno.ParamErr)
	return
}

// Upload implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Upload(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	resp := new(video.PublishListResponse)
	if req == nil {
		resp.Resp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err := service.NewVideoItemService(ctx).CreateVideoItem(req) // 创建 video item
	if err != nil {
		klog.Errorf("save video item error: %v", err.Error())
		resp.Resp = pack.BuildBaseResp(errno.PublishVideoErr)
	} else {
		resp.Resp = pack.BuildBaseResp(errno.Success)
	}
	return resp, nil
}

// GetVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideoList(ctx context.Context, req *video.GetVideoListRequest) (resp *video.
	GetVideoListResponse, err error) {
	var (
		ids        []int64
		user2Video = map[int64]*video.Video{}
		users      *userPb.UserListResponse
	)
	resp = new(video.GetVideoListResponse)
	if req == nil {
		goto ParamErr
	}
	if req.GetAuthor_ {
		if req.Author == 0 {
			goto ParamErr
		}
		resp.VideoList, err = service.NewFeedsService(ctx).GetUserPublish(req) // 用户发布的
	} else if req.Uuid == 0 {
		if req.IsLike {
			goto ParamErr
		}
		resp.VideoList, err = service.NewFeedsService(ctx).FeedsItem(req) // 未登录用户
	} else {
		// Uuid != 0 说明是登陆用户
		if req.IsLike {
			resp.VideoList, err = service.NewLikeService(ctx).LikesItem(req) // 获取喜欢列表
		} else {
			resp.VideoList, err = service.NewFeedsService(ctx).LoginUserFeedsItem(req) // 登陆用户获取, 需要 返回是否喜欢
		}
	}

	klog.Infof("get video list :%#v", resp.VideoList)
	if err != nil || len(resp.VideoList) == 0 {
		resp.Resp = pack.BuildBaseResp(errno.GetVideoErr)
		return resp, nil
	}
	ids = make([]int64, 0, len(resp.VideoList))
	for i := 0; i < len(resp.VideoList); i++ {
		uId := resp.VideoList[i].Author.Id
		ids = append(ids, uId)
		user2Video[uId] = resp.VideoList[i]
	}
	users, err = s.userRpc.GetUsers(ctx, &userPb.GetUserSRequest{Id: ids})
	if err != nil {
		klog.Errorf("获取视频用户失败 %v", err)
		resp.Resp = pack.BuildBaseResp(errno.GetVideoErr)
		return resp, nil
	}
	for i := 0; i < len(users.User); i++ {
		user2Video[users.User[i].Id].Author = users.User[i]
	}

	resp.Resp = pack.BuildBaseResp(errno.Success)
	return resp, nil

ParamErr:
	resp.Resp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

func (s *VideoServiceImpl) UpdateVideoInfoConStart() func() {
	cons := mq.NewSubConsumer(constants.VActionVideoComCountQCount, constants.VActionVideoComCountExName,
		mq.VGetActionVideoComQueueName, mq.VGetActionVideoComCountQueueKey, "")
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
					_, err = s.IncrComment(context.Background(), &video.IncrCommentRequest{
						VideoId: info.VideoId, //TODO:
						Add:     info.ActionType == 1,
					})
					if err != nil {
						klog.Errorf("%d 号 消息队列处理失败: DB保存失败%v", i, err)
						continue
					}
					break
				}
			}

		}(i, cons[i])
	}
	cleanUp := func() {
		done <- struct{}{}
	}
	return cleanUp

}

//ConsumerStart 开启消费者 监听 Save.Video. 消息队列
func (s *VideoServiceImpl) ConsumerStart() func() {
	cons := mq.NewSubConsumer(constants.VideoQCount, constants.SaveVideoExName, mq.GetSaveVideoQueueName, mq.GetSaveVideoQueueKey, "")
	factory := storage.DefaultOssFactory{
		Key: constants.OssSecretKey,
		Id:  constants.OssSecretID,
		Url: constants.OssUrl,
	}
	upload := factory.Factory()

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
				data.Body, err = util.DeCompress(data.Body)
				if err != nil {
					klog.Errorf("%d 号 消息解压失败:%v", i, err)
					return
				}
				buf := bytes.NewReader(data.Body)

				if err != nil {
					klog.Errorf("%d 号 消息接受失败:%v", i, err)
					continue
				}

				decoder := gob.NewDecoder(buf)
				if err != nil {
					//TODO:继续发送
					klog.Errorf("%d 号 消息队列处理失败:%v", i, err)
					continue

				}
				var info storage.Info
				err = decoder.Decode(&info)
				if err != nil {
					return
				}
				klog.Infof("%d 号 消息队列获取到参数:%s %d %d", i, info.Title, info.Uuid, time.Unix(info.Time,
					0).Format(constants.TimeFormatS))

				playUrl, coverUrl := upload.UploadFile(&info)
				if playUrl == "" {
					// 上传失败
					klog.Errorf("上传失败, 检查 oss 服务是否正常工作")
					continue
				}
				for i := 0; i < 2; i++ {
					_, err = s.Upload(context.Background(), &video.PublishListRequest{
						Author:   info.Uuid,
						PlayUrl:  playUrl,
						CoverUrl: coverUrl,
						Title:    info.Title,
					})
					if err != nil {
						klog.Errorf("%d 号 消息队列处理失败: DB保存失败%v", i, err)
						continue
					}
					break
				}
			}

		}(i, cons[i])
	}
	cleanUp := func() {
		done <- struct{}{}
	}

	return cleanUp
}
