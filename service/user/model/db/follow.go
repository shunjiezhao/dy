package db

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
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
	// 1.维护 follow_list 里面的 follow_each_other
	isNew := true
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var db *gorm.DB
		if err := updateUserFollowInfo(tx, followerId, id, add); err != nil {
			return err
		}
		if add {
			db = tx.Where(follow).Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil}),
			}).Create(&follow)

		} else {
			db = tx.Where(follow).Delete(&Follow{})
		}
		if err := db.Error; err != nil {
			log.Println("社交操作: 失败 ", err.Error())
			return err
		}
		isNew = db.RowsAffected != 0
		if !isNew {
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
	return isNew, err
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

//GetFollowUserList 获取关注列表
func GetFollowUserList(ctx context.Context, fromUserId int64) ([]*User, error) {
	userList, err := getFollowUserHelper(ctx, fromUserId)
	if err != nil {
		return nil, err
	}
	// 获取关注的列表, 将是否关注改为 true
	for i := 0; i < len(userList); i++ {
		userList[i].IsFollow = true
	}
	return userList, nil
}
func getFollowUserHelper(ctx context.Context, fromUserId int64) ([]*User, error) {
	var followList []*Follow
	var userList []*User
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. 查询关注列表
		db := tx.Model(&Follow{}).Order("to_user_uuid DESC").Where("from_user_uuid = ?", fromUserId).Find(&followList)
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
			userIDs[i] = followList[i].ToUserUuid
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
	// 如果是获取关注列表,那么isFollow 就是 true
	return userList, nil
}

var getFollowerSql = `
select follower.fansID 'follower_id', if(follow.to_user_uuid is not null, true, false) 'is_follow'  
from (  
         (select from_user_uuid 'fansID', to_user_uuid  
          from follow_list  
          where to_user_uuid = %d and deleted_at is null) follower -- 粉丝  
             left join  
             (select to_user_uuid -- 自己的关注  
              from follow_list  
              where from_user_uuid = %d and deleted_at is null) follow  on follow.to_user_uuid  = follower.fansID  -- 关注了自己的粉丝  
         ) order by follower_id desc;
`

//GetFollowerUserList 获取粉丝列表
func GetFollowerUserList(ctx context.Context, toUserId int64) ([]*User, error) {
	type Result struct {
		FollowerID int64 `gorm:"column:"follower_id"`
		IsFollow   bool  `gorm:"column:"is_follow"`
	}
	var (
		results  []Result
		userList []*User
	)
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Raw(fmt.Sprintf(getFollowerSql, toUserId, toUserId)).Scan(&results).Error
		if err != nil {
			return err
		}

		// 查询 该userid
		dataLen := len(results)
		userIDs := make([]int64, dataLen)
		for i := 0; i < dataLen; i++ {
			userIDs[i] = results[i].FollowerID
		}
		userList, err = MGetUsers(tx, ctx, userIDs)
		if len(userList) != dataLen {
			fmt.Printf("%#v %#v\n", results, userList)
			panic("发现错误, 一个事务内 不可重复读现象")
		}
		// 两个结果都是根据uuid 排序
		for i := 0; i < dataLen; i++ {
			userList[i].IsFollow = results[i].IsFollow
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userList, nil
}
