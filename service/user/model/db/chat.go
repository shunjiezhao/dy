package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

func SaveChat(db *gorm.DB, ctx context.Context, message *Message) error {
	return db.WithContext(ctx).Create(message).Error
}

var getChatSql = `
   select  from_user_uuid, to_user_uuid, content,created_at,message_id from message_info where (to_user_uuid, 
from_user_uuid) = (%d,
%d) or (from_user_uuid, to_user_uuid) = (%d,%d);
`

func GetChatList(db *gorm.DB, ctx context.Context, fromUserId int64, toUserId int64) (message []*Message,
	err error) {
	sql := fmt.Sprintf(getChatSql, fromUserId, toUserId, fromUserId, toUserId)
	err = db.WithContext(ctx).Raw(sql).Scan(&message).Error
	return message, err
}
