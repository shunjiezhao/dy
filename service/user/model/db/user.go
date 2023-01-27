// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package db

import (
	"context"
	"first/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Uuid          int64 `gorm:"primarykey, column:uuid"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	UserName      string         `json:"username" gorm:"column:username"`
	Password      string         `json:"password" gorm:"column:password"`
	NickName      string         `json:"nickname" gorm:"column:nickname"`
	FollowCount   int64          `json:"follow_count" gorm:"column:follow_count"`
	FollowerCount int64          `json:"follower_count" gorm:"column:follower_count"`
	IsFollow      bool           `json:"is_follow" gorm:"-"`
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// MGetUsers multiple get list of user info
func MGetUsers(db *gorm.DB, ctx context.Context, userIDs []int64) ([]*User, error) {
	if db == nil {
		db = DB.WithContext(ctx)
	}
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := db.Where("uuid in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUsers create user info
func CreateUsers(ctx context.Context, users []*User) (int64, error) {
	tx := DB.WithContext(ctx).Create(users)
	return tx.RowsAffected, tx.Error
}

func CreateUser(ctx context.Context, users *User) (int64, error) {
	tx := DB.WithContext(ctx).Create(users)
	return users.Uuid, tx.Error
}

// QueryUsers query list of user info
func QueryUsers(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryUserByName query  user info by UserName
func QueryUserByName(ctx context.Context, userName string) (*User, error) {
	var res User
	if err := DB.WithContext(ctx).Where("username = ?", userName).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// QueryUserById query  user info by UserName
func QueryUserById(ctx context.Context, id int64, followId int64) (*User, error) {
	var user *User
	var err error
	if followId == 0 {
		user, err = getUserHelper(DB.WithContext(ctx), id)
	} else {
		// 需要查询是否关注
		err = DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			user, err = getUserHelper(tx, id)
			// 关注
			follow, err := isFollowHelper(tx, followId, id)
			if follow && err == nil {
				user.IsFollow = true
			}
			return nil
		})
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}
func getUserHelper(db *gorm.DB, id int64) (*User, error) {
	var res User
	if err := db.Where("uuid = ?", id).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
