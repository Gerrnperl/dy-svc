package controller

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, "<html><head><title>simple douyin api</title></head><body></body></html>")
}
