package controller

import (
	"main/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /douyin/relation/action/ - 关系操作
// 登录用户对其他用户进行关注或取消关注。
func FollowAction(c *gin.Context) {
	userId, err := GetUserID(c, "")
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "unauthorized",
		})
		return
	}
	followedIdStr := c.Query("to_user_id")
	followedId, err := strconv.ParseInt(followedIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	actionType := c.Query("action_type")
	err = service.FollowAction(userId, followedId, actionType)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
}
