package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"main/service"
)

type CommentActionRequest struct {
	Token       string `json:"token" binding:"required"`
	VideoId     int64  `json:"video_id" binding:"required"`
	ActionType  int32  `json:"action_type" binding:"required"`
	CommentText string `json:"comment_text"`
	CommentId   int64  `json:"comment_id"`
}

type CommentResponse struct {
	Response
	Comment service.CommentInfo `json:"comment,omitempty"`
}

type CommentsResponse struct {
	Response
	CommentList []*service.CommentInfo `json:"comment_list,omitempty"`
}

// POST /douyin/comment/action/ - 评论操作
// 登录用户对视频进行评论。
func CommentAction(c *gin.Context) {
	var err error
	userId, err := GetUserID(c, "")
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请登录后再进行操作",
		})
		return
	}
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentIdStr := c.Query("comment_id")

	videoId, err1 := strconv.ParseInt(videoIdStr, 10, 64)
	commentId, err2 := strconv.ParseInt(commentIdStr, 10, 64)

	if (videoIdStr == "" || actionType == "") ||
		(actionType == "1" && commentText == "") ||
		(actionType == "2" && commentIdStr == "") ||
		(err1 != nil) ||
		(actionType == "2" && err2 != nil) ||
		(actionType != "1" && actionType != "2") {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "参数错误",
		})
		return
	}

	var comment *service.CommentInfo
	if actionType == "1" {
		comment, err = service.AddComment(userId, videoId, commentText)
	} else {
		err = service.DeleteComment(userId, commentId)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Errorf("评论操作失败: %w", err).Error(),
		})
		return
	}

	if actionType == "1" {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "评论成功",
			},
			Comment: *comment,
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "删除评论成功",
		})
	}
}

// GET /douyin/comment/list/ - 视频评论列表
// 查看视频的所有评论，按发布时间倒序。
func CommentList(c *gin.Context) {
	requestId, err := GetUserID(c, "")
	if err != nil {
		requestId = 0
	}
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if videoIdStr == "" || err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "参数错误",
		})
		return
	}

	comments, err := service.GetCommentsByVideoId(videoId, requestId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Errorf("获取评论列表失败: %w", err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CommentsResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取评论列表成功",
		},
		CommentList: comments,
	})
}
