package service

import (
	"main/models"
	"strconv"
)

// FavoriteAction 收藏/取消收藏视频
//
// creates or deletes a favorite record in the database.
// If actionType is 1, it creates a favorite record.
// If actionType is 2, it deletes a favorite record.
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

// FavoriteList 获取用户收藏列表
func FavoriteList(userId int64) ([]*models.Video, error) {
	favorites, err := models.FavoriteDao().GetVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}
