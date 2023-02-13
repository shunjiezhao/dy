package common

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

// User
//TODO: 参数检验
type (
	RegisterRequest struct {
		UserName string `json:"username" form:"username" `
		PassWord string `json:"password" form:"password"`
	}
	RegisterResponse struct {
		handlers.Response
		handlers.UserId
		handlers.Token
	}
	LoginRequest struct {
		UserName string `json:"username"  form:"username"`
		PassWord string `json:"password" form:"password"`
		UserId   int64  // 查询用户信息使用
	}

	LoginResponse struct {
		handlers.Response
		handlers.UserId
		handlers.Token
	}
	GetInfoRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetInfoResponse struct {
		handlers.Response
		*handlers.User `json:"user"`
	}

	GetUserSRequest struct {
		Id        []int64 `json:"id,omitempty" form:"id" binding:"required"`
		CurUserId int64
	}
)

func SendRegisterResponse(c *gin.Context, userId int64, token string) {
	c.JSON(consts.StatusOK, RegisterResponse{
		Response: handlers.BuildResponse(errno.Success),
		UserId:   handlers.UserId{UserId: userId},
		Token:    handlers.Token{Token: token},
	})
}
func SendGetInfoResponse(c *gin.Context, user *handlers.User) {
	c.JSON(consts.StatusOK, GetInfoResponse{
		Response: handlers.BuildResponse(errno.Success),
		User:     user,
	})
}

// Follow

type (
	FollowActionType int32 // 1-关注，2-取消关注
	ActionRequest    struct {
		handlers.Token
		FromUserId       int64 `form:"-" json:"-"`
		ToUserId         int64 `json:"to_user_id" form:"to_user_id"`
		FollowActionType `json:"action_type" form:"action_type"`
	}

	ActionResponse struct {
		handlers.Response
	}

	GetUserFollowerListRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetUserFollowerListResponse struct {
		handlers.Response
		Users []*handlers.User `json:"user_list,omitempty"`
	}

	GetUserFollowListRequest struct {
		handlers.UserId
		handlers.Token
	}
	GetUserFollowListResponse struct {
		handlers.Response
		Users []*handlers.User `json:"users,omitempty"`
	}
)

func (f FollowActionType) IsFollow() bool {
	return f == 1
}

type (
	CommentActionType int32 // 1-发布评论，2-删除评论

	// Comment

	CommentActionRequest struct {
		handlers.Token
		VideoId           int64 `json:"video_id" form:"video_id"`
		CommentActionType `json:"action_type" form:"action_type"`
		CommentText       string `json:"comment_text"form:"comment_text"`
		CommentId         int64  `json:"comment_id" form:"comment_id"`
		UserId            int64  `json:"-" form:"-"`
	}
	CommentActionResponse struct {
		handlers.Response
		*handlers.Comment `json:"comment"`
	}
	CommentListRequest struct {
		VideoId int64 `json:"video_id" form:"video_id"`
		handlers.Token
	}
	CommentListResponse struct {
		CommentList []*handlers.Comment `json:"comment_list"`
		handlers.Response
	}
)

// IsAdd 是否是发布评论
func (f CommentActionType) IsAdd() bool {
	return f == 1
}

func (f CommentActionType) String() string {
	if f.IsAdd() {
		return "添加"
	}
	return "删除"
}

type (
	// Chat
	ChatActionRequest struct {
		handlers.Token
		handlers.ToUserId
		ActionType int32  `json:"action_type" form:"action_type"` // 1-发送消息
		Content    string `json:"content" form:"content"`
	}
	ChatActionResponse struct {
		handlers.Response
		*handlers.Message
	}

	ChatListRequest struct {
		handlers.Token
		handlers.ToUserId
	}

	ChatListResponse struct {
		handlers.Response
		MessageList []*handlers.Message `json:"message_list"` // 消息列表
	}
	//好友列表
	FriendListRequest struct {
		handlers.Token
		handlers.UserId
	}

	FriendListResponse struct {
		handlers.Response
		UserList []*handlers.FriendUser `json:"user_list"` // 用户列表
	}
)

func SendCommentListResponse(c *gin.Context, comments []*handlers.Comment) {
	c.JSON(consts.StatusOK, CommentListResponse{
		Response:    handlers.BuildResponse(errno.Success),
		CommentList: comments,
	})
}

func SendUserListResponse(c *gin.Context, users []*handlers.User) {
	c.JSON(consts.StatusOK, GetUserFollowerListResponse{
		Response: handlers.BuildResponse(errno.Success),
		Users:    users,
	})
}
func SendCommentResponse(c *gin.Context, comment *handlers.Comment) {
	c.JSON(consts.StatusOK, CommentActionResponse{
		Response: handlers.BuildResponse(errno.Success),
		Comment:  comment,
	})
}
func SendChatResponse(c *gin.Context, msg *handlers.Message) {
	c.JSON(consts.StatusOK, ChatActionResponse{
		Response: handlers.BuildResponse(errno.Success),
		Message:  msg,
	})
}

func GetChatListResponse(c *gin.Context, msg []*handlers.Message) {
	c.JSON(consts.StatusOK, ChatListResponse{
		Response:    handlers.BuildResponse(errno.Success),
		MessageList: msg,
	})
}
func GetFriendListResponse(c *gin.Context, msg []*handlers.FriendUser) {
	c.JSON(consts.StatusOK, FriendListResponse{
		Response: handlers.BuildResponse(errno.Success),
		UserList: msg,
	})
}
