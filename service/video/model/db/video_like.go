package db

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateFavVideo(ctx context.Context, FavVideo *FavouriteVideo) error {
	FavVideo.IsLike = true
	err := VideoDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		db := tx.Table(constants.FavouriteVideoTableName).
			Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil})}).
			Create(FavVideo)

		err := db.Error
		if err != nil {
			klog.Errorf("[Video-DB]: 喜欢操作失败: %v", err.Error())
			return err

		}

		isNew := db.RowsAffected != 0
		if !isNew {
			klog.Errorf("[Video-DB]: 喜欢操作先前存在")
			return nil

		}

		err = likeUpdVideoInfo(tx, FavVideo.VideoId, true)
		if err != nil {
			klog.Errorf("[Video-DB]: 喜欢操作更新视频信息失败: %v", err.Error())
			return err
		}
		return nil

	})
	return err
}
func DeleteFavVideo(ctx context.Context, FavVideo *FavouriteVideo) error {
	err := VideoDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		db := tx.Table(constants.FavouriteVideoTableName).Where("uuid = ? and video_id = ?", FavVideo.Uuid,
			FavVideo.VideoId).Delete(&FavouriteVideo{})

		err := db.Error
		if err != nil {
			klog.Errorf("[Video-DB]: 取消喜欢操作失败: %v", err.Error())
			return err
		}

		if db.RowsAffected == 0 {
			klog.Errorf("[Video-DB]: 喜欢操作不存在")
			return errno.RecordNotExistErr

		}

		err = likeUpdVideoInfo(tx, FavVideo.VideoId, false)
		if err != nil {
			klog.Errorf("[Video-DB]: 取消喜欢操作更新视频信息失败: %v", err.Error())
			return err
		}
		return nil

	})
	return err
}

var op = map[bool]string{
	false: "favourite_count - ?",
	true:  "favourite_count + ?",
}

// 喜欢/ 不喜欢 对于视频信息的影响
func likeUpdVideoInfo(db *gorm.DB, videoId int64, isLike bool) error {
	// 更新视频条目
	db = db.Model(&Video{}).Where("video_id = ?", videoId).Updates(map[string]interface{}{"favourite_count": gorm.
		Expr(op[isLike], 1)})
	err := db.Error
	return err
}

// GetFavVideoAfterTime 返回 t 时间之前的 count 个 视频, 按照发布时间降序
func GetFavVideoAfterTime(ctx context.Context, uuid, t int64, count int) ([]*FavouriteVideo, error) {
	videos := make([]*FavouriteVideo, 0)

	if err := VideoDb.WithContext(ctx).Table(constants.FavouriteVideoTableName).Preload("Video").Order("created_at DESC").
		Where("created_at < ? and uuid = ?", util.GetMysqlTime(t), uuid).
		Limit(count).Find(&videos).Error; err != nil {

		return nil, err
	}
	return videos, nil
}
