package db

import (
	"context"
	"first/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Follow struct {
	gorm.Model
	FromUserUuid int64 `gorm:"column:from_user_uuid"`
	ToUserUuid   int64 `gorm:"column:to_user_uuid"`
}

func (f *Follow) TableName() string {
	return constants.FollowTableName
}

// IsFollow 查询是否 followerId 是 id 的粉丝
func IsFollow(ctx context.Context, id int64, followerId int64) (bool, error) {
	return isFollowHelper(DB.WithContext(ctx), id, followerId)
}
func isFollowHelper(db *gorm.DB, id int64, followerId int64) (bool, error) {
	var res Follow
	err := db.Where(&Follow{
		FromUserUuid: followerId,
		ToUserUuid:   id,
	}).First(&res).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FollowUser(ctx context.Context, id int64, followerId int64) (bool, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.OnConflict{DoNothing: true}).Create(&Follow{
			FromUserUuid: followerId,
			ToUserUuid:   id,
		}).Error
		if err != nil {
			//TODO:日志
			return err
		}
		err = updateUserFollowInfo(tx, followerId, id, true)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
func UnFollowUser(ctx context.Context, id int64, followerId int64) (bool, error) {
	// from 取消关注 to
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.OnConflict{DoNothing: true}).Delete(&Follow{
			FromUserUuid: followerId,
			ToUserUuid:   id,
		}).Error
		if err != nil {
			//TODO:日志
			return err
		}
		err = updateUserFollowInfo(tx, followerId, id, false)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func updateUserFollowInfo(tx *gorm.DB, followerId int64, id int64, inc bool) error {
	op := "+"
	if !inc {
		op = "-"
	}
	//增加 user 的 follow_count 和 follower_count
	if err := tx.Model(&User{}).Where("uuid = ?", followerId).Updates(map[string]interface{}{"follow_count": gorm.
		Expr(
			"follow_count "+op+" ?",
			1)}).Error; err != nil {
		//TODO:日志
		return err
	}
	if err := tx.Model(&User{}).Where("uuid = ?",
		id).Updates(map[string]interface{}{"follower_count": gorm.Expr(
		"follower_count "+op+" ?",
		1)}).Error; err != nil {
		//TODO:日志
		return err
	}
	return nil
}

//GetFollowList 获取关注列表
func GetFollowList(ctx context.Context, fromUserId int64) ([]*User, error) {
	return getFollowListByIdStat(ctx, fromUserId, "from_user_uuid = ?")
}

//GetFollowerList 获取粉丝列表
func GetFollowerList(ctx context.Context, toUserId int64) ([]*User, error) {
	return getFollowListByIdStat(ctx, toUserId, "to_user_uuid = ?")
}

func getFollowListByIdStat(ctx context.Context, id int64, query string) ([]*User, error) {
	var followList []*Follow
	var userList []*User
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. 查询关注列表
		err := tx.Where(query, id).Find(&followList).Error
		if err != nil {
			return err
		}
		var userIDs []int64
		// 将 关注的 用户的 ID 提取出来
		for _, follow := range followList {
			userIDs = append(userIDs, follow.ToUserUuid)
		}
		//2. 获取用户列表
		userList, err = MGetUsers(tx, ctx, userIDs)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return userList, nil
}
