package service

import (
	"fmt"
	"main/models"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
)

func TestUserRegister(t *testing.T) {
	username := "testuser"
	password := "testpassword"

	user := &models.User{
		Id:       1,
		Name:     username,
		Password: password,
	}

	patch := gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "Add", func(*models.UserDaoStruct, *models.User) (*models.User, error) {
		return user, nil
	})
	defer patch.Reset()

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	id, token2, err := UserRegister(username, password)
	if err != nil {
		t.Fatalf("UserRegister failed: %v", err)
	}

	if id != user.Id {
		t.Errorf("UserRegister returned wrong ID: got %d, want %d", id, user.Id)
	}

	if token2 != token {
		t.Errorf("UserRegister returned wrong token: got %s, want %s", token2, token)
	}
}

func TestUserLogin_Success(t *testing.T) {
	// Arrange
	username := "testuser"
	password := "testpassword"
	user := &models.User{
		Id:       1,
		Name:     username,
		Password: password,
	}
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	patch := gomonkey.ApplyFunc(Authenticate, func(username, password string) (int64, error) {
		return user.Id, nil
	})
	defer patch.Reset()

	// Act
	id, token2, err := UserLogin(username, password)

	// Assert
	if err != nil {
		t.Fatalf("UserLogin failed: %v", err)
	}
	if id != user.Id {
		t.Errorf("UserLogin returned wrong ID: got %d, want %d", id, user.Id)
	}
	if token2 != token {
		t.Errorf("UserLogin returned wrong token: got %s, want %s", token2, token)
	}
}

func TestUserLogin_AuthenticateError(t *testing.T) {
	// Arrange
	username := "testuser"
	password := "testpassword"
	expectedErr := models.ErrNotFound{Model: "user", Key: "name", Value: username}
	patch := gomonkey.ApplyFunc(Authenticate, func(username, password string) (int64, error) {
		return -1, expectedErr
	})
	defer patch.Reset()

	// Act
	id, token, err := UserLogin(username, password)

	// Assert
	if id != -1 {
		t.Errorf("UserLogin returned wrong ID: got %d, want -1", id)
	}
	if token != "" {
		t.Errorf("UserLogin returned wrong token: got %s, want empty string", token)
	}
	if err == nil {
		t.Fatalf("UserLogin should have returned an error")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("UserLogin returned wrong error: got %v, want %v", err, expectedErr)
	}
}

func TestUserProfile_Success(t *testing.T) {
	// Arrange
	userId := int64(1)
	rawUser := &models.User{
		Id:              userId,
		Name:            "testuser",
		FollowCount:     10,
		FollowerCount:   20,
		Avatar:          "avatar.png",
		BackgroundImage: "background.png",
		Signature:       "test signature",
		TotalFavorited:  30,
		WorkCount:       40,
		FavoriteCount:   50,
	}
	expectedUserProfile := &UserProfile{
		Id:              userId,
		Name:            rawUser.Name,
		FollowCount:     rawUser.FollowCount,
		FollowerCount:   rawUser.FollowerCount,
		Avatar:          rawUser.Avatar,
		BackgroundImage: rawUser.BackgroundImage,
		Signature:       rawUser.Signature,
		TotalFavorited:  rawUser.TotalFavorited,
		WorkCount:       rawUser.WorkCount,
		FavoriteCount:   rawUser.FavoriteCount,
	}
	patch := gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "GetById", func(*models.UserDaoStruct, int64) (*models.User, error) {
		return rawUser, nil
	})
	defer patch.Reset()

	// Act
	userProfile, err := GetUserProfile(userId, 0)

	// Assert
	if err != nil {
		t.Fatalf("UserProfile failed: %v", err)
	}
	if !reflect.DeepEqual(userProfile, expectedUserProfile) {
		t.Errorf("UserProfile returned wrong user profile: got %v, want %v", userProfile, expectedUserProfile)
	}
}

func TestUserProfile_UserDaoError(t *testing.T) {
	// Arrange
	userId := int64(1)
	expectedErr := fmt.Errorf("failed to get user by ID")
	patch := gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "GetById", func(*models.UserDaoStruct, int64) (*models.User, error) {
		return nil, expectedErr
	})
	defer patch.Reset()

	// Act
	userProfile, err := GetUserProfile(userId, 0)

	// Assert
	if userProfile != nil {
		t.Errorf("UserProfile should have returned nil user profile")
	}
	if err == nil {
		t.Fatalf("UserProfile should have returned an error")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("UserProfile returned wrong error: got %v, want %v", err, expectedErr)
	}
}
