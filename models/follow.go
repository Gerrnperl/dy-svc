package models

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model

	FollowerId int64 `gorm:"index" json:"follow_id"`
	FollowedId int64 `gorm:"index" json:"followed_id"`
}

func (f *Follow) TableName() string {
	return "follow"
}

var (
	_followDaoInstance *FollowDaoStruct
	_followDaoOnce     sync.Once
)

type FollowDaoStruct struct{}

func FollowDao() *FollowDaoStruct {
	_followDaoOnce.Do(func() {
		_followDaoInstance = &FollowDaoStruct{}
	})
	return _followDaoInstance
}

func (dao *FollowDaoStruct) FollowAction(follow *Follow, do bool) error {
	if do {
		// Check if the follow already exists
		var count int64
		if err := DB().Model(&Follow{}).Where("follower_id = ? AND followed_id = ?", follow.FollowerId, follow.FollowedId).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("follow relation already exists")
		}
		if err := DB().Create(follow).Error; err != nil {
			return err
		}
	} else {
		if err := DB().Unscoped().Where("follower_id = ? AND followed_id = ?", follow.FollowerId, follow.FollowedId).Delete(&Follow{}).Error; err != nil {
			return err
		}
	}

	delta := 1

	if !do {
		delta = -1
	}

	// Update the follower's follow count
	if err := DB().Model(&User{}).Where("id = ?", follow.FollowerId).Update("follow_count", gorm.Expr("follow_count + ?", delta)).Error; err != nil {
		return err
	}

	// Update the followed user's follower count
	if err := DB().Model(&User{}).Where("id = ?", follow.FollowedId).Update("follower_count", gorm.Expr("follower_count + ?", delta)).Error; err != nil {
		return err
	}

	return nil
}

func (dao *FollowDaoStruct) IsFollowing(followerId int64, followedId int64) (bool, error) {
	var count int64
	if err := DB().Model(&Follow{}).Where("follower_id = ? AND followed_id = ?", followerId, followedId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
