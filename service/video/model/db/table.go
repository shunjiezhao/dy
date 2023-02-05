package db

import (
	"first/pkg/constants"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt time.Time      `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"column:video_id"`
	AuthorUuid    int64  `json:"author" gorm:"column:author_uuid"` // author uuid
	Title         string `json:"title" gorm:"title"`
	PlayUrl       string `json:"play_url,omitempty" gorm:"column:play_url"`
	CoverUrl      string `json:"cover_url,omitempty" gorm:"column:cover_url"`
	FavoriteCount int64  `json:"favorite_count,omitempty" gorm:"column:favourite_count"`
	CommentCount  int64  `json:"comment_count,omitempty" gorm:"column:comment_count"`

	Base
}

type FavouriteVideo struct {
	Uuid    int64 `json:"author" gorm:"column:uuid"`
	VideoId int64
	Video   Video `gorm:"foreignKey:video_id;references:VideoId;"` // belong to 属于
	IsLike  bool  `json:"is_like" gorm:"is_like"`

	Base
}

func (*FavouriteVideo) TableName() string {
	return constants.UserTableName
}
func (u *Video) TableName() string {
	return constants.VideoTableName
}
