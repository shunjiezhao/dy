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

package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode                = 0
	ServiceErrCode             = 10001
	ParamErrCode               = 10002
	UserAlreadyExistErrCode    = 10003
	AuthorizationFailedErrCode = 10004
	UserAlreadyFollowErrCode   = 10005
	RecordNotExistErrCode      = 10006
	RecordAlreadyExistErrCode  = 10007
	RemoteErrCode              = 10008
	PublishVideoErrCode        = 10009
	GetVideoErrCode            = 10010
	LikeOpErrCode              = 10011
	OpSelfErrCode              = 10012
	MsgSaveErrCode             = 10013
	CompressErrCode            = 10014
	UnCompressErrCode          = 10015
	RedisKeyNotExistErrCode    = 10016
	VideoSaveErrCode           = 10017
	VideoBrokeErrCode          = 10018
	RemoteOssErrCode           = 10019
)

type ErrNo struct {
	ErrCode int64  `json:"status_code"`
	ErrMsg  string `json:"status_msg"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success                = NewErrNo(SuccessCode, "Success")
	ServiceErr             = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamErr               = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	UserAlreadyExistErr    = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	UserAlreadyFollowErr   = NewErrNo(UserAlreadyFollowErrCode, "User already follow")
	AuthorizationFailedErr = NewErrNo(AuthorizationFailedErrCode, "用户验证失败")
	RecordNotExistErr      = NewErrNo(RecordNotExistErrCode, "record not exist")
	RecordAlreadyExistErr  = NewErrNo(RecordAlreadyExistErrCode, "record already exist")
	RemoteErr              = NewErrNo(RemoteErrCode, "稍后重试")
	PublishVideoErr        = NewErrNo(PublishVideoErrCode, "上传视频失败")
	GetVideoErr            = NewErrNo(GetVideoErrCode, "获取视频失败")
	LikeOpErr              = NewErrNo(LikeOpErrCode, "喜欢操作失败")
	OpSelfErr              = NewErrNo(OpSelfErrCode, "不能对自己操作哦")
	MsgSaveErr             = NewErrNo(MsgSaveErrCode, "发送消息失败")
	CompressErr            = NewErrNo(CompressErrCode, "压缩失败")
	UnCompressErr          = NewErrNo(UnCompressErrCode, "解压压缩失败")
	RedisKeyNotExistErr    = NewErrNo(RedisKeyNotExistErrCode, "Redis Key 不存在")
	VideoSaveErr           = NewErrNo(VideoSaveErrCode, "Video 保存失败")
	VideoBrokeErr          = NewErrNo(VideoBrokeErrCode, "视频 损坏失败")
	RemoteOssErr           = NewErrNo(RemoteOssErrCode, "稍后重试")
)

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}

	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
