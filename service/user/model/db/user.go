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
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// MGetUsers multiple get list of user info
func MGetUsers(ctx context.Context, userIDs []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUser create user info
func CreateUser(ctx context.Context, users []*User) (int64, error) {
	tx := DB.WithContext(ctx).Create(users)
	return tx.RowsAffected, tx.Error
}

// QueryUsers query list of user info
func QueryUsers(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryUser query  user info by UserName
func QueryUser(ctx context.Context, userName string) (*User, error) {
	var res User
	if err := DB.WithContext(ctx).Where("username = ?", userName).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
