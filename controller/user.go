package controller

import (
	"main/models"

	"github.com/gin-gonic/gin"
)

type UserRegisterResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

// POST /douyin/user/register/ - 用户注册接口
// 新用户注册时提供用户名，密码，昵称即可，用户名需要保证唯一。创建成功后返回用户 id 和权限token.
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, err := models.AddUser(&models.User{
		Name:     username,
		Password: password,
	})

	if err != nil {
		c.JSON(200, UserRegisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	token, err := models.GenerateToken(user)

	if err != nil {
		c.JSON(200, UserRegisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(200, UserRegisterResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: user.Id,
		Token:  token,
	})
}
