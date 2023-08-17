package controller

import (
	"main/models"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFavoriteListResponse struct {
	Response
	VideoList []models.Video `json:"video_list"`
}

// POST /douyin/favorite/action/ - 赞操作
// 登录用户对视频的点赞和取消点赞操作。
func FavoriteAction(c *gin.Context) {
	userId, err := GetUserID(c, "")
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "unauthorized",
		})
		return
	}
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	err = service.FavoriteAction(userId, videoId, actionType)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
}

// GET /douyin/favorite/list/ - 喜欢列表
// 登录用户的所有点赞视频。
func FavoriteList(c *gin.Context) {
	userId, err := GetUserID(c, c.Query("user_id"))
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "unauthorized",
		})
		return
	}
	list, err := service.FavoriteList(userId)
	service.AdjustVideosUrl(list)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	videoList := make([]models.Video, len(list))
	for i, v := range list {
		videoList[i] = *v
	}
	c.JSON(http.StatusOK, GetFavoriteListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
