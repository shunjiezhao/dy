package db

import (
	"context"
	videoPb "first/kitex_gen/video"
	"first/pkg/constants"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 主要测试 喜欢操作 以及是否正确获取到视频信息
func insertVideoHelper(n int) []*Video {
	var video []*Video
	for i := 0; i < n; i++ {
		video = append(video, &Video{
			Id:         int64(i + 1),
			AuthorUuid: 999,
			Title:      "电影-1",
			PlayUrl:    "url",
			CoverUrl:   "url",
			Base: Base{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}
	return video
}

func TestFavourite(t *testing.T) {
	//1.创建数据
	InitVideo()
	DB, err := VideoDb.DB()
	if err != nil {
		t.Fatalf("DB 初始化失败")
	}
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	db := VideoDb
	if db == nil {
		t.Fatalf("DB 初始化失败")
	}

	assert := assert.New(t)
	ctx := context.Background()

	//1. 创建 3 个 video
	nVideos := 3
	videos := insertVideoHelper(3)
	for _, video := range videos {
		db.Unscoped().Delete(video, video.Id)
	}
	if err := db.CreateInBatches(videos, len(videos)).Error; err != nil { // 插入
		t.Fatalf("DB 插入失败")
	}

	// 喜欢关系 uuid->video_id
	input := "1->1 1->2 2->3 2->2"
	favs := getFavHelper(input, true)
	for _, fav := range favs {
		db.Unscoped().Table(constants.FavouriteVideoTableName).Where("uuid = ?  ", fav.Uuid).Delete(&FavouriteVideo{})
	}
	for i := 0; i < len(favs); i++ {

		err := CreateFavVideo(ctx, favs[i])
		if err != nil {
			t.Fatalf("喜欢关系 插入失败")
		}
	}
	now := time.Now().Add(time.Hour).Unix()
	packFav := func(uuid int64) ([]*videoPb.Video, error) {
		time.Sleep(time.Millisecond * 100)
		item, err := GetFavVideoAfterTime(ctx, uuid, now, nVideos)
		if err != nil {
			return nil, err

		}
		return FavVideos(item), nil
	}

	// 1. 喜欢喜欢的
	// 2. 喜欢不喜欢的
	// 3. 取消喜欢
	// 4. 再次取消喜欢
	// 5. 给用户 1 返回视频列表, 查看喜欢的情况是否正确
	// 4. 返回用户 2 的列表, 查看是否正确

	// 检查 正确
	// 1. 喜欢的视频信息相同
	// 2. 喜欢数相同
	tests := []struct {
		name              string
		op                func() ([]*videoPb.Video, error)
		wantLike          map[int]bool
		wantLen           int64
		wantVideoFavCount map[int]int
		shouldErr         bool
	}{
		{
			name: "喜欢喜欢的",
			op: func() ([]*videoPb.Video, error) {
				err := CreateFavVideo(ctx, &FavouriteVideo{Uuid: 1, VideoId: 1})
				if err != nil {
					return nil, err

				}
				return packFav(1)
			},
			wantLen:           2,
			wantLike:          map[int]bool{1: true, 2: true, 3: false},
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 1},
		},
		{
			name: "喜欢不喜欢的",
			op: func() ([]*videoPb.Video, error) {
				err := CreateFavVideo(ctx, &FavouriteVideo{Uuid: 1, VideoId: 3})

				if err != nil {
					return nil, err

				}
				return packFav(1)
			},
			wantLen:           3,
			wantLike:          map[int]bool{1: true, 2: true, 3: true},
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 2},
		},
		{
			name: "取消喜欢",
			op: func() ([]*videoPb.Video, error) {
				err := DeleteFavVideo(ctx, &FavouriteVideo{Uuid: 1, VideoId: 3})
				if err != nil {
					return nil, err

				}

				return packFav(1)
			},
			wantLen:           2,
			wantLike:          map[int]bool{1: true, 2: true, 3: false},
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 1},
		},
		{
			name: "再次取消喜欢",
			op: func() ([]*videoPb.Video, error) {
				return nil, DeleteFavVideo(ctx, &FavouriteVideo{Uuid: 1, VideoId: 3})
			},
			shouldErr:         true,
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 1},
		},
		{
			name: "给用户 1 返回视频列表, 查看喜欢的情况是否正确",
			op: func() ([]*videoPb.Video, error) {
				item, err := LoginUserFeedsItem(ctx, now, 1)
				if err != nil {
					return nil, err
				}

				return LoginFeeds(item), nil
			},
			wantLen:           -1,
			wantLike:          map[int]bool{1: true, 2: true, 3: false},
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 1},
		},
		{
			name: "返回用户 2 的列表, 查看是否正确",
			op: func() ([]*videoPb.Video, error) {
				item, err := LoginUserFeedsItem(ctx, now, 2)
				if err != nil {
					return nil, err
				}

				return LoginFeeds(item), nil
			},
			wantLen:           -1,
			wantLike:          map[int]bool{2: true, 1: false, 3: true},
			wantVideoFavCount: map[int]int{1: 1, 2: 2, 3: 1},
		},
	}

	for _, test := range tests {

		items, err := test.op()
		if test.shouldErr {
			assert.NotNil(err, "%s: ", test.name)
			continue
		}
		if test.wantLen != -1 {
			assert.Equal(int(test.wantLen), len(items), "%s: ", test.name)
		}
		for _, item := range items {
			if item.Id > 0 && item.Id < 4 {
				assert.Equal(test.wantLike[int(item.Id)], item.IsFavorite, "%s: ", test.name)
				assert.Equal(test.wantVideoFavCount[int(item.Id)], int(item.FavoriteCount), "%s: ", test.name)
				assert.Equal(videos[int(item.Id)-1].Title, item.Title, "%s: ", test.name)
			}
		}
	}

}
func getFavHelper(input string, like bool) []*FavouriteVideo {
	split := strings.Split(input, " ")
	likes := make([]*FavouriteVideo, len(split))
	for i, s := range split {
		ids := strings.Split(s, "->")
		uuid, _ := strconv.ParseInt(ids[0], 10, 64)
		videoId, _ := strconv.ParseInt(ids[1], 10, 64)
		likes[i] = &FavouriteVideo{
			Uuid:    uuid,
			VideoId: videoId,
			IsLike:  like,
		}

	}
	return likes
}
func FavVideos(dVideos []*FavouriteVideo) []*videoPb.Video {
	videos := make([]*videoPb.Video, len(dVideos))
	for i := 0; i < len(dVideos); i++ {
		videos[i] = video(&dVideos[i].Video)
		videos[i].IsFavorite = true // 获取喜欢的短视频
	}
	return videos
}

func video(dVideo *Video) *videoPb.Video {
	return &videoPb.Video{
		Id:            dVideo.Id,
		Author:        dVideo.AuthorUuid,
		PlayUrl:       dVideo.PlayUrl,
		CoverUrl:      dVideo.CoverUrl,
		FavoriteCount: dVideo.FavoriteCount,
		CommentCount:  dVideo.CommentCount,
		Title:         dVideo.Title,
	}
}

func LoginFeeds(dVideos []*LoginFeedResult) []*videoPb.Video {
	res := make([]*videoPb.Video, len(dVideos))
	for i := 0; i < len(res); i++ {
		res[i] = video(&dVideos[i].Video)
		res[i].IsFavorite = dVideos[i].IsFavourite
	}
	return res
}
