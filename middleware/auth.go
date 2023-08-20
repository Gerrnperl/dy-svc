package middleware

import (
	"main/controller"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// auth
//
// authenticates the user token and sets the user ID in the context.
// It returns the user ID if authentication is successful, or an error otherwise.
func auth(c *gin.Context, token string) (int64, error) {
	id, err := service.AuthenticateToken(token)
	if err != nil {
		return 0, err
	}
	if id > 0 {
		c.Set("user_id", id)
	}
	return id, nil
}

// AuthQuery
//
// a middleware that authenticates the user token from the query string.
// It calls the auth function to authenticate the token and set the user ID in the context.
func AuthQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		_, err := auth(c, token)
		if err != nil {
			c.Set("auth_err", err)
		}
		c.Next()
	}
}

func AuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		_, err := auth(c, token)
		if err != nil {
			c.Set("auth_err", err)
		}
		c.Next()
	}
}

func AuthBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		_, err := auth(c, token)
		if err != nil {
			c.Set("auth_err", err)
		}
		c.Next()
	}
}

// PassAuth
//
// aborts the request if the user is not authenticated.
func PassAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, _ := c.Get("auth_err")
		if err != nil {
			c.JSON(http.StatusUnauthorized, controller.Response{
				StatusCode: 1,
				StatusMsg:  err.(error).Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
