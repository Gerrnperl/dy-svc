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

	apiRouter.GET("/feed/", controller.Feed)

	apiRouter.POST("/user/register/", controller.UserRegister)

	apiRouter.POST("/user/login/", controller.UserLogin)

	apiRouter.GET("/user/", middleware.AuthQuery(), controller.UserProfile)

	apiRouter.POST("/publish/action/", middleware.AuthBody(), controller.UploadVideo)

	apiRouter.GET("/publish/list/", middleware.AuthQuery(), controller.GetPublishList)

	apiRouter.POST("/favorite/action/", middleware.AuthQuery(), controller.FavoriteAction)

	apiRouter.GET("/favorite/list/", middleware.AuthQuery(), controller.FavoriteList)

	apiRouter.POST("/comment/action/", middleware.AuthQuery(), controller.CommentAction)

	apiRouter.GET("/comment/list", middleware.AuthQuery(), controller.CommentList)
}
