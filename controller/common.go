package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

// GetUserID
//
// extracts the user ID from the given *gin.Context object.
//
//	@param c *gin.Context: The *gin.Context object to extract the user ID from.
//	@param userId string: The user ID as a string. If this parameter is not empty, the function tries to parse it as an int64 value and returns it. If this parameter is empty, the function tries to get the user ID from the context using c.Get("user_id").
func GetUserID(c *gin.Context, userId string) (int64, error) {
	var id int64

	if userId == "" {
		userId, existed := c.Get("user_id")
		var ok bool
		id, ok = userId.(int64)
		if !ok || !existed {
			return 0, errors.New("user id is not given")
		}
	} else {
		var err error
		id, err = strconv.ParseInt(userId, 10, 64)
		if err != nil {
			return 0, errors.New("user id is not valid")
		}
	}

	return id, nil
}
