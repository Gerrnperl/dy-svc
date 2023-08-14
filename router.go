package main

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	apiRouter.GET("/hello/", controller.Hello)

	apiRouter.POST("/user/register/", controller.UserRegister)
}
