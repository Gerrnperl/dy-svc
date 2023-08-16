package controller

import (
	"main/models"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserCredentialsResponse - 用户注册/登录接口返回的数据结构
type UserCredentialsResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"` // 用户 id
	Token  string `json:"token,omitempty"`   // 用户 token
}

type UserProfilesResponse struct {
	Response
	User *models.UserProfile `json:"user,omitempty"`
}

// POST /douyin/user/register/ - 用户注册接口
// 新用户注册时提供用户名，密码，昵称即可，用户名需要保证唯一。创建成功后返回用户 id 和权限token.
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.UserRegister(username, password)

	if err != nil {
		c.JSON(200, UserCredentialsResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(200, UserCredentialsResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: id,
		Token:  token,
	})
}

// POST /douyin/user/login/ - 用户登录接口
// 通过用户名和密码进行登录，登录成功后返回用户 id 和权限 token.
func UserLogin(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	id, token, err := service.UserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(200, UserCredentialsResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: id,
		Token:  token,
	})

}

// GET /douyin/user/ - 用户信息
// 获取登录用户的 id、昵称，如果实现社交部分的功能，还会返回关注数和粉丝数。
func UserProfile(c *gin.Context) {

	userId, err := GetUserID(c, c.Query("userId"))

	if err != nil {
		c.JSON(200, UserProfilesResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	user, err := service.UserProfile(userId)

	if err != nil {
		c.JSON(200, UserProfilesResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(200, UserProfilesResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		User: user,
	})
}
