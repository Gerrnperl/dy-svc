package service

import "main/models"

type CommentInfo struct {
	Id         int64       `json:"id,omitempty"`
	User       UserProfile `json:"user,omitempty"`
	Content    string      `json:"content,omitempty"`
	CreateDate string      `json:"create_date,omitempty"` // "mm-dd"
}

// AddComment 添加评论
//
// creates a new comment record in the database and returns the comment info.
func AddComment(userId, videoId int64, commentText string) (comment *CommentInfo, err error) {
	rawComment := models.Comment{
		UserId:  userId,
		VideoId: videoId,
		Content: commentText,
	}
	user, err := GetUserProfile(userId, 0)
	if err != nil {
		return nil, err
	}
	err = models.CommentDao().CreateComment(&rawComment)
	if err != nil {
		return nil, err
	}
	comment = &CommentInfo{
		Id:         rawComment.Id,
		User:       *user,
		Content:    rawComment.Content,
		CreateDate: rawComment.CreatedAt.Format("01-02"),
	}
	return comment, nil
}

// DeleteComment 删除评论
func DeleteComment(userId, commentId int64) error {
	return models.CommentDao().DeleteComment(userId, commentId)
}

// GetCommentsByVideoId 根据视频id获取评论
func GetCommentsByVideoId(videoId int64, requestId int64) ([]*CommentInfo, error) {
	rawComments, err := models.CommentDao().GetCommentsByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	comments := make([]*CommentInfo, len(rawComments))
	for i, rawComment := range rawComments {
		user, err := GetUserProfile(rawComment.UserId, requestId)
		if err != nil {
			return nil, err
		}
		comments[i] = &CommentInfo{
			Id:         rawComment.Id,
			User:       *user,
			Content:    rawComment.Content,
			CreateDate: rawComment.CreatedAt.Format("01-02"),
		}
	}
	return comments, nil
}
