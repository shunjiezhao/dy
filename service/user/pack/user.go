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

package pack

import (
	"errors"
	user "first/kitex_gen/user"
	"first/pkg/errno"
	"first/service/user/model/db"
)

// BuildBaseResp build baseResp from error
func BuildBaseResp(err error) *user.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}
func baseResp(err errno.ErrNo) *user.BaseResp {
	return &user.BaseResp{StatusCode: err.ErrCode, StatusMsg: err.ErrMsg}
}

func Users(dUsers []*db.User) []*user.User {
	users := make([]*user.User, 0)
	for i := 0; i < len(dUsers); i++ {
		users = append(users, User(dUsers[i]))
	}
	return users
}
func User(dUser *db.User) *user.User {
	return &user.User{
		Id:            dUser.Uuid,
		UserName:      dUser.UserName,
		FollowCount:   dUser.FollowCount,
		FollowerCount: dUser.FollowerCount,
		IsFollow:      dUser.IsFollow,
	}
}
