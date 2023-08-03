package dao

import (
	"fmt"
)

// SelectMessageChat 查询聊天记录
func SelectMessageChat(userId int64, toUserId int64) ([]*Message, error) {
	isFriend, err := IsFriend(userId, int64(toUserId))
	if err != nil {
		return nil, fmt.Errorf("select friend relationship error:%v", err)
	}
	if !isFriend {
		return nil, fmt.Errorf("frined no exist")
	}
	var chats = make([]*Message, 0, 100)
	res := db.Where("from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ? ", userId, toUserId, toUserId, userId).Find(&chats)
	if res.Error != nil {
		return nil, fmt.Errorf("select messagechat error:%v", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return chats, nil
}

// IsFriend 是否为好友
func IsFriend(userId, toUserId int64) (bool, error) {
	var friend Follow
	db := Connection()
	res := db.Select("id").Where("following_user_id = ? AND followed_user_id = ? AND relationship = ?", userId, toUserId, 1).Find(&friend)
	if res.Error != nil {
		return false, fmt.Errorf("select freind relaotionship error:%v", res.Error)
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
