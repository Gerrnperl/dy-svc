package models

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	Id int64 `json:"id,omitempty" gorm:"primarykey"`

	UserId  int64 `json:"user_id,omitempty"`
	VideoId int64 `json:"video_id,omitempty"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

var (
	_favoriteDaoInstance *FavoriteDaoStruct
	_favoriteDaoOnce     sync.Once
)

type FavoriteDaoStruct struct{}

func FavoriteDao() *FavoriteDaoStruct {
	_favoriteDaoOnce.Do(func() {
		_favoriteDaoInstance = &FavoriteDaoStruct{}
	})
	return _favoriteDaoInstance
}

// Action
//
// adds or removes a favorite.
// It will also update the FavoriteCount of the video,
// the TotalFavorited of the author and the FavoriteCount of the user.
// If do is true, it adds the favorite; otherwise, it removes the favorite.
// If the favorite already exists, it will be removed. (seems that the demo app does not support "unfavorite")
func (d *FavoriteDaoStruct) Action(f *Favorite, do bool) error {
	var err error
	if do {
		// check if the favorite exists
		var count int64
		err = DB().Model(&Favorite{}).Where("user_id = ? AND video_id = ?", f.UserId, f.VideoId).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return d.Action(f, false)
		}
		err = DB().Create(&f).Error
	} else {
		// hard delete
		err = DB().Unscoped().Where("user_id = ? AND video_id = ?", f.UserId, f.VideoId).Delete(&Favorite{}).Error
	}
	if err != nil {
		return err
	}
	delta := 1
	if !do {
		delta = -1
	}
	// Update the FavoriteCount of the video.
	err = DB().Model(&Video{}).Where("id = ?", f.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error
	if err != nil {
		return err
	}
	// Update the TotalFavorited of the author.
	// get the author id of the video
	var video Video
	err = DB().Where("id = ?", f.VideoId).First(&video).Error
	if err != nil {
		return err
	}
	// Update User's TotalFavorited
	err = DB().Model(&User{}).Where("id = ?", video.AuthorId).Update("total_favorited", gorm.Expr("total_favorited + ?", delta)).Error
	if err != nil {
		return err
	}
	// Update the FavoriteCount of the user
	err = DB().Model(&User{}).Where("id = ?", f.UserId).Update("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error
	if err != nil {
		return err
	}
	return nil
}

// (soft) delete ALL favorites of a video.
// designed to be called when a video is deleted.
func (d *FavoriteDaoStruct) DeleteByVideoId(videoId int64) error {
	// soft delete
	err := DB().Where("video_id = ?", videoId).Delete(&Favorite{}).Error
	return err
}

func (d *FavoriteDaoStruct) GetByUserId(userId int64) ([]*Favorite, error) {
	var favorites []*Favorite
	err := DB().Where("user_id = ?", userId).Find(&favorites).Error
	return favorites, err
}

func (d *FavoriteDaoStruct) GetByVideoId(videoId int64) ([]*Favorite, error) {
	var favorites []*Favorite
	err := DB().Where("video_id = ?", videoId).Find(&favorites).Error
	return favorites, err
}

func (d *FavoriteDaoStruct) GetUsersByVideoId(videoId int64) ([]*User, error) {
	var users []*User
	// Query the database to find all users who have favorited a specific video.
	// The query uses a left join to combine the favorite and user tables, and selects all columns from the user table.
	err := DB().
		Table("favorite").
		Select("user.*").
		Joins("left join user on user.id = favorite.user_id").
		Where("favorite.video_id = ?", videoId).
		Find(&users).
		Error
	return users, err
}

func (d *FavoriteDaoStruct) GetVideosByUserId(userId int64) ([]*Video, error) {
	var videos []*Video
	err := DB().
		Table("favorite").
		Select("video.*").
		Joins("join video on video.id = favorite.video_id").
		Where("favorite.user_id = ?", userId).
		Find(&videos).
		Error
	return videos, err
}

func (d *FavoriteDaoStruct) GetUsersCountByVideoId(videoId int64) (int64, error) {
	var count int64
	err := DB().
		Model(&Favorite{}).
		Where("video_id = ?", videoId).
		Count(&count).
		Error
	return count, err
}

func (d *FavoriteDaoStruct) GetVideosCountByUserId(userId int64) (int64, error) {
	var count int64
	err := DB().
		Model(&Favorite{}).
		Where("user_id = ?", userId).
		Count(&count).
		Error
	return count, err
}
