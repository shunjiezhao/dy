package db

import (
	"context"
	"first/pkg/constants"
	"first/pkg/util"
)

func CreateCommentS(ctx context.Context, items []*Comment) (int64, error) {
	tx := DB.WithContext(ctx).Create(items)
	return tx.RowsAffected, tx.Error
}

func CreateComment(ctx context.Context, Comment *Comment) (*Comment, error) {
	tx := DB.WithContext(ctx).Create(Comment)
	return Comment, tx.Error
}
func DeleteComment(ctx context.Context, commentId int64) error {
	tx := DB.WithContext(ctx).Delete(&Comment{}, commentId)
	return tx.Error
}

// GetCommentAfterTime 返回 t 时间之前的 count 个 评论, 按照发布时间降序
func GetCommentAfterTime(ctx context.Context, videoId int64, t int64, count int) ([]*Comment, error) {
	comments := make([]*Comment, 0)

	if err := DB.WithContext(ctx).Table(constants.CommentTableName).Preload("User").Order("created_at DESC").
		Where(
			"video_id = ? and created_at <= ?", videoId,
			util.GetMysqlTime(t)).Limit(count).Find(&comments).Error; err != nil {

		return nil, err
	}
	return comments, nil
}

// GetComment 返回 t 时间之前的 count 个 评论, 按照发布时间降序
func GetComment(ctx context.Context, videoId int64, count int) ([]*Comment, error) {
	comments := make([]*Comment, 0)

	if err := DB.WithContext(ctx).Table(constants.CommentTableName).Preload("User").Order("created_at DESC").
		Where("video_id = ? ", videoId).Limit(count).Find(&comments).Error; err != nil {

		return nil, err
	}
	return comments, nil
}
