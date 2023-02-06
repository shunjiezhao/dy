package pack

import (
	"first/kitex_gen/video"
	"first/service/api/handlers"
)

func Videos(videos []*video.Video) []*handlers.Video {
	res := make([]*handlers.Video, len(videos))
	for i := 0; i < len(videos); i++ {

		res[i] = &handlers.Video{
			Id:            videos[i].Id,
			PlayUrl:       videos[i].PlayUrl,
			CoverUrl:      videos[i].CoverUrl,
			FavoriteCount: videos[i].FavoriteCount,
			CommentCount:  videos[i].CommentCount,
			IsFavorite:    videos[i].IsFavorite,
			Author:        &handlers.User{Id: videos[i].Author},
		}

	}
	return res

}
