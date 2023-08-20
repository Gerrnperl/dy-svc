package controller

import (
	"fmt"
	"main/models"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPublishListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

// POST /douyin/publish/action/ - 视频投稿
// 登录用户选择视频上传。
func UploadVideo(c *gin.Context) {
	userIdStr, existed := c.Get("user_id")
	if !existed {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "用户未登录",
		})
	}
	userId := userIdStr.(int64)

	data, err := c.FormFile("data")

	title := c.PostForm("title")

	if err != nil {
		c.JSON(400, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Errorf("上传文件失败: %v", err).Error(),
		})
		return
	}

	if title == "" {
		c.JSON(400, Response{
			StatusCode: 1,
			StatusMsg:  "标题不能为空",
		})
		return
	}

	_, err = service.UploadVideo(userId, data, title)

	if err != nil {
		c.JSON(400, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Errorf("上传文件失败: %v", err).Error(),
		})
		return
	}

	c.JSON(200, Response{
		StatusCode: 0,
		StatusMsg:  "上传成功",
	})
}

// GET /douyin/publish/list/ - 发布列表
// 登录用户的视频发布列表，直接列出用户所有投稿过的视频。
func GetPublishList(c *gin.Context) {
	userId, err := GetUserID(c, c.Query("user_id"))

	if err != nil {
		c.JSON(200, UserProfilesResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	publishList, err := service.GetPublishList(userId)

	if err != nil {
		c.JSON(400, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Errorf("获取发布列表失败: %v", err).Error(),
		})
		return
	}

	videoList := make([]models.Video, len(publishList))
	for i, v := range publishList {
		videoList[i] = *v
	}

	c.JSON(200, GetPublishListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
