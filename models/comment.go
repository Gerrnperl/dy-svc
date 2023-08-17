package models

import (
	"sync"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Id      int64  `json:"id,omitempty" gorm:"primarykey"`
	VideoId int64  `json:"video_id,omitempty"`
	UserId  int64  `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}

func (c *Comment) TableName() string {
	return "comment"
}

type CommentDaoStruct struct{}

var (
	_commentDaoInstance *CommentDaoStruct
	_commentDaoOnce     sync.Once
)

func CommentDao() *CommentDaoStruct {
	_commentDaoOnce.Do(func() {
		_commentDaoInstance = &CommentDaoStruct{}
	})
	return _commentDaoInstance
}

func (dao *CommentDaoStruct) CreateComment(comment *Comment) error {
	err := DB().Create(&comment).Error
	if err != nil {
		return err
	}
	// add video comment count
	err = DB().Model(&Video{}).Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	return err
}

func (dao *CommentDaoStruct) GetCommentById(id int64) (*Comment, error) {
	comment := &Comment{}
	err := DB().Where("id = ?", id).First(comment).Error
	return comment, err
}

func (dao *CommentDaoStruct) GetCommentsByVideoId(videoId int64) ([]*Comment, error) {
	comments := []*Comment{}
	err := DB().Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (dao *CommentDaoStruct) DeleteComment(userId, commentId int64) error {
	err := DB().Where("id = ? and user_id = ?", commentId, userId).Delete(&Comment{}).Error
	if err != nil {
		return err
	}
	// minus video comment count
	err = DB().Model(&Video{}).Where("id = ?", commentId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	return err
}
