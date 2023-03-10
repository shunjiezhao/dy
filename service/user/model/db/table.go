package db

import (
	"first/pkg/constants"
	"gorm.io/gorm"
	"time"
)

type (
	Base struct {
		CreatedAt time.Time      `json:"created_at" gorm:"created_at"`
		UpdatedAt time.Time      `json:"updated_at" gorm:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
	//Comment 评论信息表
	Comment struct {
		Id      int64 `json:"id,omitempty" gorm:"column:comment_id"`
		Uuid    int64
		User    User   `gorm:"foreignKey:uuid;references:Uuid;"`
		VideoId int64  `json:"video_id"  gorm:"column:video_id"`
		Content string `json:"content" gorm:"column:content"`
		Base
	}
	//User 用户信息表
	User struct {
		Uuid int64 `gorm:"primarykey, column:uuid"`
		Base

		UserName      string `json:"username" gorm:"column:username"`
		Password      string `json:"password" gorm:"column:password"`
		NickName      string `json:"nickname" gorm:"column:nickname"`
		FollowCount   int64  `json:"follow_count" gorm:"column:follow_count"`
		FollowerCount int64  `json:"follower_count" gorm:"column:follower_count"`
		IsFollow      bool   `json:"is_follow" gorm:"is_follow"`
	}
	//Follow 关注信息表
	Follow struct {
		Base

		FromUserUuid int64 `gorm:"column:from_user_uuid"`
		ToUserUuid   int64 `gorm:"column:to_user_uuid"`
	}

	//Message 评论信息表
	Message struct {
		Id           int64  `gorm:"column:message_id"`
		FromUserUuid int64  `gorm:"column:from_user_uuid"`
		ToUserUuid   int64  `gorm:"column:to_user_uuid"`
		Content      string `gorm:"column:content"`

		Base
	}
)

func (f *Follow) TableName() string {
	return constants.FollowTableName
}
func (u *User) TableName() string {
	return constants.UserTableName
}
func (*Comment) TableName() string {
	return constants.CommentTableName
}
func (*Message) TableName() string {
	return constants.MessageTableName
}
