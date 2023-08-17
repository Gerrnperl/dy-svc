package service

import (
	"main/models"
	"strconv"
)

func FavoriteAction(userId int64, videoId string, actionType string) error {
	vid, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		return err
	}
	action, err := strconv.ParseInt(actionType, 10, 32)
	if err != nil {
		return err
	}
	return models.FavoriteDao().Action(&models.Favorite{
		UserId:  userId,
		VideoId: vid,
	}, action == 1)
}

func FavoriteList(userId int64) ([]*models.Video, error) {
	favorites, err := models.FavoriteDao().GetVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}
