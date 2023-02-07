package constants

var (
	// MQ
	MQConnURL = "amqp://rabbit:123456@localhost:5672/"

	// SAVE Video
	SaveVideoExName       = "视频保存"
	SaveVideoPrefix       = "dy.save.video."
	SaveVideoKey          = "saveVideo"
	VideoQCount     int64 = 4

	//  Add Comment
	AddCommentExName = "添加评论"
	AddCommentPrefix = "dy.add.comment."
)
