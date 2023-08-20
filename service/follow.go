package service

import "main/models"

func FollowAction(followerId int64, followedId int64, actionType string) error {
	return models.FollowDao().FollowAction(&models.Follow{
		FollowerId: followerId,
		FollowedId: followedId,
	}, actionType == "1")
}
