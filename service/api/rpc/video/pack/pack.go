package pack

import (
	"first/kitex_gen/video"
	"first/service/api/handlers"
	pack2 "first/service/api/rpc/user/pack"
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
			Author:        pack2.User(videos[i].Author),
		}

	}
	return res

}
