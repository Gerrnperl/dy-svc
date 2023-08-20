package service

import (
	"errors"
	"main/models"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddCommentWithMock(t *testing.T) {
	mockUser := &UserProfile{
		Id:   1,
		Name: "testuser",
	}

	mockComment := &models.Comment{
		Id:      1,
		UserId:  1,
		VideoId: 1,
		Content: "test comment",
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	}

	patch1 := gomonkey.ApplyFunc(GetUserProfile, func(userId int64) (*UserProfile, error) {
		return mockUser, nil
	})
	defer patch1.Reset()

	patch2 := gomonkey.ApplyMethod(reflect.TypeOf(models.CommentDao()), "CreateComment", func(dao *models.CommentDaoStruct, comment *models.Comment) error {
		return nil
	})
	defer patch2.Reset()

	comment, err := AddComment(1, 1, "test comment")

	assert.NoError(t, err)

	expectedComment := &CommentInfo{
		User:    *mockUser,
		Content: mockComment.Content,
	}
	assert.Equal(t, expectedComment.User, comment.User)
	assert.Equal(t, expectedComment.Content, comment.Content)
}

func TestAddCommentUserProfileErrorWithMock(t *testing.T) {
	patch := gomonkey.ApplyFunc(GetUserProfile, func(userId int64) (*UserProfile, error) {
		return nil, errors.New("error getting user profile")
	})
	defer patch.Reset()

	comment, err := AddComment(0, 1, "test comment")

	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestAddCommentCreateCommentErrorWithMock(t *testing.T) {
	mockUser := &UserProfile{
		Id:   1,
		Name: "testuser",
	}

	patch1 := gomonkey.ApplyFunc(GetUserProfile, func(userId int64) (*UserProfile, error) {
		return mockUser, nil
	})
	defer patch1.Reset()

	patch2 := gomonkey.ApplyMethod(reflect.TypeOf(models.CommentDao()), "CreateComment", func(dao *models.CommentDaoStruct, comment *models.Comment) error {
		return errors.New("error creating comment")
	})
	defer patch2.Reset()

	comment, err := AddComment(1, 1, "test comment")

	assert.Error(t, err)
	assert.Nil(t, comment)
}
