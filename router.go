package main

import (
	"EdgeTB-backend/handler"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/EdgeTB")

	//用户 apis
	{
		apiRouter.POST("/login", handler.Login)
		apiRouter.POST("/register", handler.Register)
		apiRouter.POST("/logout", handler.Logout)
		apiRouter.POST("/getUsername", handler.GetUsername)
		apiRouter.GET("/getUserInfo", handler.GetUserInfo)
		apiRouter.POST("/setUserInfo", handler.UpdateUserInfo)
		apiRouter.POST("/setPassword", handler.UpdateUserPwd)
	}
}
