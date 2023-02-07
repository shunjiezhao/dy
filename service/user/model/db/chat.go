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

var getFriendChatListSql = `
select other_id, mySend as 'my_send' , content
from(

        select other_id, row_number() over (partition by other_id order by created_at desc ) as rk,mySend,content
        from (
                 select to_user_uuid 'other_id', true as 'mySend', created_at , 
row_number() over (partition by from_user_uuid, to_user_uuid order by created_at desc ) as rk , content from message_info where from_user_uuid = %d
                 union all
                 select from_user_uuid 'other_id', false as 'mySend', created_at, row_number() over (partition by  from_user_uuid,to_user_uuid order by created_at desc ) as rk , content from message_info where to_user_uuid = %d
             )
                 t where rk = 1
    ) tmp where rk = 1;

`

type FriendChatListResult struct {
	Content string `gorm:"content"`
	OtherId int64  `gorm:"other_id"`
	MySend  bool   `gorm:"my_send"`
}

func GetFriendChatList(db *gorm.DB, ctx context.Context, fromUserId int64) (message []*FriendChatListResult,
	err error) {
	sql := fmt.Sprintf(getFriendChatListSql, fromUserId, fromUserId)
	err = db.WithContext(ctx).Raw(sql).Scan(&message).Error
	if err != nil {
		return nil, err
	}

	return message, err
}
