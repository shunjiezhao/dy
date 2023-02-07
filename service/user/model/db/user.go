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
	"first/pkg/errno"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

var pool sync.Pool

func init() {
	pool = sync.Pool{New: func() any {
		return &strings.Builder{}
	}}
}

// MGetUsers multiple get list of user info order by uuid DESC
func MGetUsers(db *gorm.DB, ctx context.Context, userIDs []int64) ([]*User, error) {
	if db == nil {
		db = DB.WithContext(ctx)
	}
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := db.Order("uuid DESC").Where("uuid in ?", userIDs).Find(&res).Error; err != nil {
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

// QueryUsersById query list of user info
func QueryUsersById(ctx context.Context, uuid []int64) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("uuid in ?", uuid).Find(&res).Error; err != nil {
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

// QueryUserByNamePwd query  user info by UserName
func QueryUserByNamePwd(ctx context.Context, userName, passWord string) (*User, error) {
	var res User
	if t := DB.WithContext(ctx).Where("username = ? and password = ?", userName,
		passWord).First(&res).RowsAffected; t == 0 {
		return nil, errno.RecordNotExistErr
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

var getUserSLoginSql = `
select  l.uuid as 'uuid', l.username as 'username',  l.follow_count as 'follow_count', 
l.follower_count as 'follower_count', l.nickname as 'nickname',
      if(r.to_user_uuid is not null or uuid = %d,true,false)  'is_follow'
from (
	  select uuid, username,  follow_count, follower_count, nickname
	  from user_info where  uuid in (%s)
) l 
left join 
(select to_user_uuid from follow_list where from_user_uuid = %d) r 
on  l.uuid = r.to_user_uuid;
`

func GetUserSLogin(ctx context.Context, id []int64, uuid int64) ([]*User, error) {
	if len(id) == 0 {
		return nil, nil
	}

	builder := pool.Get().(*strings.Builder)
	builder.WriteString(strconv.FormatInt(id[0], 10))
	for i := 1; i < len(id); i++ {
		builder.WriteString("," + strconv.FormatInt(id[i], 10)) // 1, 2, 3
		// 避免判断
	}
	// 1, 2, 3,
	klog.Infof("得到id %s", builder.String())

	var res []*User
	sql := fmt.Sprintf(getUserSLoginSql, uuid, builder.String(), uuid)
	err := DB.WithContext(ctx).Raw(sql).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil

}
