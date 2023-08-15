package main

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	r.GET("/", controller.Hello)

	apiRouter.POST("/user/register/", controller.UserRegister)

	apiRouter.POST("/user/login/", controller.UserLogin)

	apiRouter.GET("/user/", middleware.AuthQuery(), controller.UserProfile)
}
