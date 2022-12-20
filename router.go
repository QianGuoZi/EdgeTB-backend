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

	//数据集 apis
	{
		//apiRouter.GET("/public-datasets", handler.AllPublicDatasets)
		//apiRouter.GET("/dataset/public/:id", handler.AllPublicDatasets)
		//apiRouter.GET("/dataset/my", handler.AllPublicDatasets)
		//apiRouter.POST("/dataset/my", handler.AllPublicDatasets)
		//apiRouter.GET("/dataset/my/:id", handler.AllPublicDatasets)
		//apiRouter.PUT("/dataset/my/:id", handler.AllPublicDatasets)
		//apiRouter.DELETE("/dataset/my/:id", handler.AllPublicDatasets)
		apiRouter.POST("/dataset/my/upload", handler.UploadDataset)
	}
}
