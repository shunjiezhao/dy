package db

import (
	"context"
	"first/pkg/constants"
	"first/pkg/logger"
	"first/pkg/util"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateVideoItem 创建 视频信息
func CreateVideoItem(db *gorm.DB, ctx context.Context, video *Video) error {
	if err := db.WithContext(ctx).Create(video).Error; err != nil {
		logger.GetLogger().Error("[DB]: 保存视频信息失败", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func IncrVideoCommentCount(db *gorm.DB, ctx context.Context, videoId int64, add bool) error {
	var op = "comment_count + ?"
	if !add {
		op = "comment_count - ?"
	}
	if err := db.WithContext(ctx).Table(constants.VideoTableName).Update("comment_count", gorm.Expr(op, 1)).Error; err != nil {
		logger.GetLogger().Error("[DB]: 保存视频信息失败", zap.String("err", err.Error()))
		return err
	}
	return nil
}

// GetUserPublish 获取用户的发布列表
func GetUserPublish(db *gorm.DB, ctx context.Context, uuid int64) ([]*Video, error) {
	if db == nil {
		db = VideoDb.WithContext(ctx)
	}
	res := make([]*Video, 0)
	if err := db.WithContext(ctx).Table(constants.VideoTableName).Order("created_at DESC").Where("author_uuid = ?",
		uuid).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetVideosAfterTime(ctx context.Context, t int64, count int) ([]*Video, error) {
	videos := make([]*Video, 0)
	// select * from video_info where create_time > t order by (created_at,create_time) DESC limit 50

	if err := VideoDb.WithContext(ctx).Table(constants.VideoTableName).Order("created_at DESC").Where("created_at < ?",
		util.GetMysqlTime(t)).Limit(count).Find(&videos).Error; err != nil {

		return nil, err
	}
	return videos, nil
}

var loginUserFeedsItemSql = `
select l.video_id,
       play_url,
       author_uuid,
       cover_url,
       title,
       favourite_count,
       comment_count,
       created_at,
       if(r.video_id is null, false, true) as 'is_favourite'
from (
         (select video_id,
                 play_url,
                 author_uuid,
                 cover_url,
                 title,
                 favourite_count,
                 comment_count,
                 created_at
          from video_info
          where created_at < '%s' and deleted_at is null order by created_at limit 50) l
             left join
             (select video_id from user_favourite_video where uuid = %d and deleted_at is null and is_like = true) r
         on l.video_id = r.video_id
         );
`

type LoginFeedResult struct {
	Video
	IsFavourite bool `gorm:"column:is_favourite"`
}

func LoginUserFeedsItem(ctx context.Context, t, uuid int64) ([]*LoginFeedResult, error) {
	var results []*LoginFeedResult
	sql := fmt.Sprintf(loginUserFeedsItemSql, util.GetMysqlTime(t), uuid)
	err := VideoDb.WithContext(ctx).Raw(sql).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
