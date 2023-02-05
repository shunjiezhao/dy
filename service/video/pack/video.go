package pack

import (
	"errors"
	videoPb "first/kitex_gen/video"
	"first/pkg/errno"
	"first/service/video/model/db"
)

// BuildBaseResp build baseResp from error
func BuildBaseResp(err error) *videoPb.VideoBaseResp {
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

func baseResp(err errno.ErrNo) *videoPb.VideoBaseResp {
	return &videoPb.VideoBaseResp{StatusCode: err.ErrCode, StatusMsg: err.ErrMsg}
}

func FavVideos(dVideos []*db.FavouriteVideo) []*videoPb.Video {
	videos := make([]*videoPb.Video, len(dVideos))
	for i := 0; i < len(dVideos); i++ {
		videos[i] = video(&dVideos[i].Video)
		videos[i].IsFavorite = true // 获取喜欢的短视频
	}
	return videos
}
func Videos(dVideos []*db.Video) []*videoPb.Video {
	videos := make([]*videoPb.Video, len(dVideos))
	for i := 0; i < len(dVideos); i++ {
		videos[i] = video(dVideos[i])
	}
	return videos
}

func video(dVideo *db.Video) *videoPb.Video {
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

func LoginFeeds(dVideos []*db.LoginFeedResult) []*videoPb.Video {
	res := make([]*videoPb.Video, len(dVideos))
	for i := 0; i < len(res); i++ {
		res[i] = video(&dVideos[i].Video)
		res[i].IsFavorite = dVideos[i].IsFavourite
	}
	return res
}
