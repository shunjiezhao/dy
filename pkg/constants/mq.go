package constants

const (
	// MQ
	MQConnURL = "amqp://rabbit:123456@localhost:5672/"

	// SAVE Video
	SaveVideoExName       = "视频保存"
	SaveVideoPrefix       = "dy.save.video."
	SaveVideoKey          = "saveVideo"
	VideoQCount     int64 = 4

	//  Action Comment
	UActionCommentExName       = "评论操作"
	UActionCommentPrefix       = "dy.action.comment."
	UActionCommentKey          = "actionComment"
	UActionCommentQCount int64 = 1

	// Update Video Info
	VActionVideoComCountExName       = "评论操作"
	VActionVideoComCountPrefix       = "action.video.comment.count."
	VActionVideoComCountKey          = "actionVideoComment"
	VActionVideoComCountQCount int64 = 1
)
