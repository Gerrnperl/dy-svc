package models

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	ToUserId   int64  `gorm:"index" json:"to_user_id"`
	FromUserId int64  `gorm:"index" json:"from_user_id"`
	Content    string `json:"content"`
}

func (m *Message) TableName() string {
	return "message"
}

var (
	_messageDaoInstance *MessageDaoStruct
	_messageDaoOnce     sync.Once
)

type MessageDaoStruct struct{}

func MessageDao() *MessageDaoStruct {
	_messageDaoOnce.Do(func() {
		_messageDaoInstance = &MessageDaoStruct{}
	})
	return _messageDaoInstance
}

// Add 添加消息
func (*MessageDaoStruct) Add(message *Message) (*Message, error) {
	if message.ToUserId == 0 {
		return nil, ErrMissingRequiredField{"to_user_id"}
	}
	if message.FromUserId == 0 {
		return nil, ErrMissingRequiredField{"from_user_id"}
	}
	if message.Content == "" {
		return nil, ErrMissingRequiredField{"content"}
	}
	// 精确到秒，防止轮询时重复
	message.CreatedAt = time.Now().Truncate(time.Second)
	if err := DB().Create(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// GetListByUserId 获取两个用户之间的消息列表
func (*MessageDaoStruct) GetListByUserId(user1, user2 int64, after time.Time) ([]*Message, error) {
	var messages []*Message
	afterStr := after.Format("2006-01-02 15:04:05")
	if err := DB().
		Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?) AND created_at > ?", user1, user2, user2, user1, afterStr).
		Order("created_at DESC").
		Find(&messages).
		Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// GetLatestConversations
//
// retrieves the latest conversations for a given user ID.
// It returns a slice of Message objects representing the latest messages in each conversation,
// sorted by creation date in descending order.
func (*MessageDaoStruct) GetLatestConversations(userId int64) ([]*Message, error) {
	var messages []*Message

	// "SELECT * FROM message
	//		WHERE id IN (
	//			SELECT MAX(id) FROM message
	//				WHERE to_user_id = ? OR from_user_id = ?
	//				GROUP BY LEAST(to_user_id, from_user_id), GREATEST(to_user_id, from_user_id)
	//		) ORDER BY created_at DESC
	// ", userId, userId

	subQuery := DB().Table("message").
		Select("MAX(id)").
		Where("to_user_id = ? OR from_user_id = ?", userId, userId).
		Group("LEAST(to_user_id, from_user_id), GREATEST(to_user_id, from_user_id)")

	if err := DB().Table("message").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	return messages, nil
}
