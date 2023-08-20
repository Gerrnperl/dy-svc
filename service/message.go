package service

import (
	"main/models"
	"time"
)

type FriendUser struct {
	UserProfile
	Message     string `json:"message"` // 最新一条消息
	MessageType int64  `json:"msgType"` // 最新一条消息类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}

type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

// PostMessage 发送信息
func PostMessage(toUserId, fromUserId int64, content string) error {
	_, err := models.MessageDao().Add(&models.Message{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
	})
	return err
}

// GetMessages 获取两个用户之间的消息列表
func GetMessages(user1Id, user2Id int64, after int64) ([]Message, error) {
	afterTime := time.Unix(after, 0)
	msgs, err := models.MessageDao().GetListByUserId(user1Id, user2Id, afterTime)
	if err != nil {
		return nil, err
	}
	var messages []Message
	for _, msg := range msgs {
		messages = append(messages, Message{
			Id:         int64(msg.ID),
			ToUserId:   msg.ToUserId,
			FromUserId: msg.FromUserId,
			Content:    msg.Content,
			CreateTime: msg.CreatedAt.Unix(),
		})
	}
	return messages, nil
}

// GetFriends 获取好友列表
//
// Friends are the users who have chatted with the current user.
// The latest message between the current user and the friend is displayed.
func GetFriends(userId int64) ([]*FriendUser, error) {
	lastMsgs, err := models.MessageDao().GetLatestConversations(userId)
	if err != nil {
		return nil, err
	}
	var friendUsers []*FriendUser
	for _, message := range lastMsgs {
		other := message.ToUserId
		messageType := int64(1)
		if other == userId {
			other = message.FromUserId
			messageType = 0
		}
		user, err := GetUserProfile(other, userId)
		if err != nil {
			return nil, err
		}
		friendUsers = append(friendUsers, &FriendUser{
			UserProfile: *user,
			Message:     message.Content,
			MessageType: messageType,
		})
	}
	return friendUsers, nil
}
