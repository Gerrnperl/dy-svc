package service

import (
	"fmt"
	"main/models"
)

type UserProfile struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	Signature       string `json:"signature,omitempty"`
}

// UserRegister
//
// registers a new user with the given username and password,
// adds the user to the database, and generates a JWT token for the user.
// Returns the user ID and token if successful, or -1 and an empty string if there is an error.
func UserRegister(username, password string) (id int64, token string, err error) {
	user, err := models.UserDao().Add(&models.User{
		Name:     username,
		Password: password,
	})
	if err != nil {
		return -1, "", err
	}

	token, err = GenerateToken(user)
	if err != nil {
		return -1, "", err
	}

	return user.Id, token, nil
}

// UserLogin
//
// authenticates a user with the given username and password,
// generates a JWT token for the user, and returns the user ID and token if successful.
// If there is an error, it returns -1 for the user ID and an empty string for the token.
func UserLogin(username, password string) (id int64, token string, err error) {
	id, err = Authenticate(username, password)
	if err != nil {
		return -1, "", err
	}

	token, err = GenerateToken(&models.User{
		Id:   id,
		Name: username,
	})
	if err != nil {
		return -1, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return id, token, nil
}

// GetUserProfile
//
// returns the user profile for the user with the given ID.
// It removes some sensitive information from the user profile, like the password.
func GetUserProfile(userId int64, requestId int64) (user *UserProfile, err error) {
	rawUser, err := models.UserDao().GetById(userId)
	if err != nil {
		return nil, err
	}
	isFollow := false
	if requestId != 0 {
		isFollow, err = models.FollowDao().IsFollowing(requestId, userId)
		if err != nil {
			return nil, err
		}
	}

	return &UserProfile{
		Id:              rawUser.Id,
		Name:            rawUser.Name,
		FollowCount:     rawUser.FollowCount,
		FollowerCount:   rawUser.FollowerCount,
		IsFollow:        isFollow,
		Avatar:          rawUser.Avatar,
		BackgroundImage: rawUser.BackgroundImage,
		Signature:       rawUser.Signature,
		TotalFavorited:  rawUser.TotalFavorited,
		WorkCount:       rawUser.WorkCount,
		FavoriteCount:   rawUser.FavoriteCount,
	}, nil
}
