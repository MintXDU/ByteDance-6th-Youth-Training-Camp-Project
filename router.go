package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.Authorize(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.Authorize(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.Authorize(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.Authorize(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.Authorize(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.Authorize(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.Authorize(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.Authorize(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.Authorize(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.Authorize(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", middleware.Authorize(), controller.FriendList)
	apiRouter.GET("/message/chat/", middleware.Authorize(), controller.MessageChat)
	apiRouter.POST("/message/action/", middleware.Authorize(), controller.MessageAction)
}
