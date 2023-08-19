package service

import (
	"main/models"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteActionAddFavoriteWithMock(t *testing.T) {
	patch := gomonkey.ApplyMethod(reflect.TypeOf(models.FavoriteDao()), "Action", func(dao *models.FavoriteDaoStruct, favorite *models.Favorite, add bool) error {
		return nil
	})
	defer patch.Reset()

	err := FavoriteAction(1, "123", "1")

	assert.NoError(t, err)
}

func TestFavoriteActionRemoveFavoriteWithMock(t *testing.T) {
	patch := gomonkey.ApplyMethod(reflect.TypeOf(models.FavoriteDao()), "Action", func(dao *models.FavoriteDaoStruct, favorite *models.Favorite, add bool) error {
		return nil
	})
	defer patch.Reset()

	err := FavoriteAction(1, "123", "0")

	assert.NoError(t, err)
}

func TestFavoriteActionInvalidVideoIdWithMock(t *testing.T) {
	err := FavoriteAction(1, "invalid", "1")

	assert.Error(t, err)
}

func TestFavoriteActionInvalidActionTypeWithMock(t *testing.T) {
	err := FavoriteAction(1, "123", "invalid")

	assert.Error(t, err)
}
