package controller

import (
	"main/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowListResponse struct {
	Response
	UserList []*service.UserProfile `json:"user_list"`
}

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

// GET /douyin/relation/follow/list/ - 用户关注列表
// 登录用户关注的所有用户列表。
func FollowList(c *gin.Context) {
	userId, err := GetUserID(c, c.Query("user_id"))
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "unauthorized",
		})
		return
	}
	followList, err := service.GetFollowings(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followList,
	})
}

// GET /douyin/relation/follower/list/ - 用户粉丝列表
// 所有关注登录用户的粉丝列表。
func FollowerList(c *gin.Context) {
	userId, err := GetUserID(c, c.Query("user_id"))
	if err != nil || userId == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  "unauthorized",
		})
		return
	}
	followerList, err := service.GetFollowers(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followerList,
	})
}
