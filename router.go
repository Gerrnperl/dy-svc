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

	apiRouter.GET("/feed/", middleware.AuthQuery(), controller.Feed)

	apiRouter.POST("/user/register/", controller.UserRegister)

	apiRouter.POST("/user/login/", controller.UserLogin)

	apiRouter.GET("/user/", middleware.AuthQuery(), middleware.PassAuth(), controller.UserProfile)

	apiRouter.POST("/publish/action/", middleware.AuthBody(), middleware.PassAuth(), controller.UploadVideo)

	apiRouter.GET("/publish/list/", middleware.AuthQuery(), middleware.PassAuth(), controller.GetPublishList)

	apiRouter.POST("/favorite/action/", middleware.AuthQuery(), middleware.PassAuth(), controller.FavoriteAction)

	apiRouter.GET("/favorite/list/", middleware.AuthQuery(), middleware.PassAuth(), controller.FavoriteList)

	apiRouter.POST("/comment/action/", middleware.AuthQuery(), middleware.PassAuth(), controller.CommentAction)

	apiRouter.GET("/comment/list", middleware.AuthQuery(), middleware.PassAuth(), controller.CommentList)

	apiRouter.POST("/relation/action/", middleware.AuthQuery(), middleware.PassAuth(), controller.FollowAction)
}
