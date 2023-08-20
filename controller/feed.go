package controller

import (
	"main/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.VideoInfo `json:"video_list"`
	NextTime  int64               `json:"next_time"`
}

// GET /douyin/Feed/ - 视频流接口
// 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个。
func Feed(c *gin.Context) {
	requestId, err := GetUserID(c, "")
	if err != nil {
		requestId = 0
	}
	// Parse request body
	var req struct {
		LatestTime int64  `form:"latest_time"`
		Token      string `form:"token"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "Bad Request",
		})
		return
	}

	if req.LatestTime == 0 {
		req.LatestTime = time.Now().Unix()
	}

	// Get videos from database
	videos, oldest, err := service.GetVideosBefore(req.LatestTime, requestId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoList := make([]service.VideoInfo, len(videos))
	for i, v := range videos {
		videoList[i] = *v
	}

	c.JSON(200, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		VideoList: videoList,
		NextTime:  oldest,
	})
}
