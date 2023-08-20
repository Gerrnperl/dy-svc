package service

import "main/models"

func FollowAction(followerId int64, followedId int64, actionType string) error {
	return models.FollowDao().FollowAction(&models.Follow{
		FollowerId: followerId,
		FollowedId: followedId,
	}, actionType == "1")
}

func GetFollowers(userId int64) ([]*UserProfile, error) {
	followers, err := models.FollowDao().GetByFollowedId(userId)
	if err != nil {
		return nil, err
	}
	var users []*UserProfile
	for _, follower := range followers {
		user, err := GetUserProfile(follower.FollowerId, userId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetFollowings(userId int64) ([]*UserProfile, error) {
	followings, err := models.FollowDao().GetByFollowerId(userId)
	if err != nil {
		return nil, err
	}
	var users []*UserProfile
	for _, following := range followings {
		user, err := GetUserProfile(following.FollowedId, userId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
