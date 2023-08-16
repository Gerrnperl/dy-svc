package service

import (
	"fmt"
	"main/models"
)

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

// UserProfile
//
// returns the user profile for the user with the given ID.
// It removes some sensitive information from the user profile, like the password.
func UserProfile(userId int64) (user *models.UserProfile, err error) {
	rawUser, err := models.UserDao().GetById(userId)
	if err != nil {
		return nil, err
	}

	return &models.UserProfile{
		Id:              rawUser.Id,
		Name:            rawUser.Name,
		FollowCount:     rawUser.FollowCount,
		FollowerCount:   rawUser.FollowerCount,
		Avatar:          rawUser.Avatar,
		BackgroundImage: rawUser.BackgroundImage,
		Signature:       rawUser.Signature,
		TotalFavorited:  rawUser.TotalFavorited,
		WorkCount:       rawUser.WorkCount,
		FavoriteCount:   rawUser.FavoriteCount,
	}, nil
}
