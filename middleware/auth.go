package middleware

import (
	"main/controller"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func auth(c *gin.Context, token string) (int64, error) {
	id, err := service.AuthenticateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, controller.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		c.Abort()
		return 0, err
	}
	if id > 0 {
		c.Set("user_id", id)
	}
	return id, nil
}

func AuthQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		_, err := auth(c, token)
		if err != nil {
			return
		}
		c.Next()
	}
}

func AuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		_, err := auth(c, token)
		if err != nil {
			return
		}
		c.Next()
	}
}

func AuthBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		_, err := auth(c, token)
		if err != nil {
			return
		}
		c.Next()
	}
}
