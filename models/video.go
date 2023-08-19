package models

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model

	Id            int64  `json:"id,omitempty" gorm:"primarykey"`
	AuthorId      int64  `json:"author_id,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	Title         string `json:"title,omitempty"`
}

func (v *Video) TableName() string {
	return "video"
}

var (
	_videoDaoInstance *VideoDaoStruct
	_videoDaoOnce     sync.Once
)

type VideoDaoStruct struct{}

func VideoDao() *VideoDaoStruct {
	_videoDaoOnce.Do(func() {
		_videoDaoInstance = &VideoDaoStruct{}
	})
	return _videoDaoInstance
}

// Add 添加视频
//
// create a new video record in the database.
// and also adds the work count of the author.
func (*VideoDaoStruct) Add(video *Video) (*Video, error) {
	if video.PlayUrl == "" {
		return nil, ErrMissingRequiredField{"play_url"}
	}
	if video.Title == "" {
		return nil, ErrMissingRequiredField{"title"}
	}
	if video.AuthorId == 0 {
		return nil, ErrMissingRequiredField{"author_id"}
	}
	// increase author's WorkCount
	if err := DB().Model(&User{}).Where("id = ?", video.AuthorId).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
		return nil, err
	}
	if err := DB().Create(&video).Error; err != nil {
		return nil, err
	}
	return video, nil
}

// GetByAuthorId 根据作者id获取视频
func (*VideoDaoStruct) GetByAuthorId(authorId int64) ([]*Video, error) {
	var videos []*Video
	if err := DB().Where("author_id = ?", authorId).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetBefore 根据时间戳获取视频
//
// It returns a list of videos created before the given timestamp.
// The number of videos returned is limited by the limit parameter.
// The oldest timestamp of the returned videos is returned as the second return value.
func (*VideoDaoStruct) GetBefore(timeStamp int64, limit int) (videoList []*Video, oldest int64, err error) {
	var videos []*Video
	// convert time to String
	timeStr := time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
	if err := DB().Where("created_at < ?", timeStr).Order("created_at desc").Limit(limit).Find(&videos).Error; err != nil {
		return nil, 0, err
	}
	if len(videos) == 0 {
		return nil, 0, nil
	}
	return videos, videos[len(videos)-1].CreatedAt.Unix(), nil
}
