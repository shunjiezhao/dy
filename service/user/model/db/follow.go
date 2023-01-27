package db

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"gorm.io/gorm"
	"log"
	"sort"
	"time"
)

type Follow struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	FromUserUuid int64          `gorm:"column:from_user_uuid"`
	ToUserUuid   int64          `gorm:"column:to_user_uuid"`
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
	// from 关注 to
	return followHelper(ctx, id, followerId, true, false)
}
func UnFollowUser(ctx context.Context, id int64, followerId int64) (bool, error) {
	// from 取消关注 to 记录需要以前存在吗?
	return followHelper(ctx, id, followerId, false, true)
}

// followHelper bool:返回之前是否存在该条记录
func followHelper(ctx context.Context, id int64, followerId int64, add bool, shouldExist bool) (bool, error) {
	follow := &Follow{
		FromUserUuid: followerId,
		ToUserUuid:   id,
	}
	// TODO: fastpath  add or delete 的情况下 查询记录是否已经存在
	isExist := true
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var db *gorm.DB
		if err := updateUserFollowInfo(tx, followerId, id, add); err != nil {
			return err
		}
		if add {
			db = tx.Where(follow).FirstOrCreate(&follow)
		} else {
			db = tx.Where(follow).Delete(&Follow{})
		}
		if err := db.Error; err != nil {
			log.Println("社交操作: 失败 ", err.Error())
			return err
		}
		isExist = db.RowsAffected != 0
		if !isExist {
			if shouldExist {
				log.Println("社交操作: 关注记录先前不存在")
				return errno.RecordNotExistErr
			}
			log.Println("社交操作: 关注记录先前存在")
			return errno.RecordAlreadyExistErr
		}

		return nil
	})
	// NOTICE: 如果已经关注, 那么直接返回成功
	if err == errno.RecordAlreadyExistErr {
		err = nil
	}
	return isExist, err
}

//updateUserFollowInfo 更新相关用户的社交信息, 需要包装在一个事务里面.
func updateUserFollowInfo(tx *gorm.DB, followerId int64, id int64, inc bool) error {
	op := "+"
	if !inc {
		op = "-"
	}
	//增加 user 的 follow_count 和 follower_count
	db := tx.Model(&User{}).Where("uuid = ?", followerId).Updates(map[string]interface{}{"follow_count": gorm.
		Expr(
			"follow_count "+op+" ?",
			1)})
	err := db.Error
	if err != nil {
		//TODO:日志
		return err
	}
	if db.RowsAffected == 0 {
		return errno.RecordNotExistErr
	}

	db = tx.Model(&User{}).Where("uuid = ?",
		id).Updates(map[string]interface{}{"follower_count": gorm.Expr(
		"follower_count "+op+" ?",
		1)})
	err = db.Error
	if err != nil {
		//TODO:日志
		return err
	}
	if db.RowsAffected == 0 {
		return errno.RecordNotExistErr
	}
	return nil
}

//TODO: 使用 join 来优化 查询关注/粉丝列表, 更新其是否关注字段

//GetFollowUserList 获取关注列表
func GetFollowUserList(ctx context.Context, fromUserId int64) ([]*User, error) {
	userList, err := getFollowListByIdStat(ctx, fromUserId, "from_user_uuid = ?", true)
	if err != nil {
		return nil, err
	}
	// 获取关注的列表, 将是否关注改为 true
	for i := 0; i < len(userList); i++ {
		userList[i].IsFollow = true
	}
	return userList, nil
}

//GetFollowerUserList 获取粉丝列表
func GetFollowerUserList(ctx context.Context, toUserId int64) ([]*User, error) {
	userList, err := getFollowListByIdStat(ctx, toUserId, "to_user_uuid = ?", false)
	// 查询是否关注自己
	if err != nil {
		return nil, err
	}
	// 获取关注的列表, 更改是否关注字段

	follows := make([]*Follow, 0)
	userIDs := make([]int64, len(userList))
	for i := 0; i < len(userList); i++ {
		userIDs[i] = userList[i].Uuid
	}
	if err := DB.WithContext(ctx).Select("to_user_uuid").Where("from_user_uuid = ? and to_user_uuid in ?", toUserId,
		userIDs).Find(&follows).Error; err != nil {
		log.Println("获取粉丝列表时 查询是否关注出错")
		return nil, err
	}
	// 更改字段
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].Uuid < userList[j].Uuid
	})
	sort.Slice(follows, func(i, j int) bool {
		return follows[i].ToUserUuid < follows[j].ToUserUuid
	})

	i, j := 0, 0
	for i != len(userList) && j != len(follows) {
		if userList[i].Uuid == follows[j].ToUserUuid { // 如果是关注的人的话
			userList[i].IsFollow = true
			i, j = i+1, j+1
		} else if userList[i].Uuid < follows[j].ToUserUuid {
			i++
		} else {
			j++
		}
	}
	return userList, nil
}

func getFollowListByIdStat(ctx context.Context, id int64, query string, isFollow bool) ([]*User, error) {
	var followList []*Follow
	var userList []*User
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. 查询关注列表
		db := tx.Where(query, id).Find(&followList)
		err := db.Error
		if err != nil {
			return err
		}
		// 如果没有关注的人
		if db.RowsAffected == 0 {
			return nil
		}
		userIDs := make([]int64, len(followList))

		// 将 关注的 用户的 ID 提取出来
		for i := 0; i < len(followList); i++ {
			if isFollow {
				userIDs[i] = followList[i].ToUserUuid
			} else {
				userIDs[i] = followList[i].FromUserUuid
			}
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
