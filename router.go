package main

import (
	"EdgeTB-backend/handler"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/EdgeTB")

	//用户 apis
	{
		apiRouter.POST("/user/login", handler.Login)
		apiRouter.POST("/user/register", handler.Register)
		apiRouter.POST("/user/logout", handler.Logout)

		apiRouter.POST("/getUsername", handler.GetUsername)
		apiRouter.GET("/getUserInfo", handler.GetUserInfo)
		apiRouter.POST("/setUserInfo", handler.UpdateUserInfo)
		apiRouter.POST("/setPassword", handler.UpdateUserPwd)
	}

	//数据集 apis
	{
		apiRouter.GET("/dataset/public", handler.AllPublicDatasets)
		apiRouter.GET("/dataset/public/:id", handler.PublicDatasetsDetail)
		apiRouter.GET("/dataset/my", handler.AllPrivateDatasets)
		apiRouter.POST("/dataset/my", handler.AddDataset)
		apiRouter.GET("/dataset/my/:id", handler.PrivateDatasetsDetail)
		apiRouter.PUT("/dataset/my/:id", handler.UpdateDataset)
		apiRouter.DELETE("/dataset/my/:id", handler.DeletePrivateDataset)
		apiRouter.POST("/dataset/my/upload", handler.UploadDataset)
	}

	//角色 apis
	{
		apiRouter.POST("/role", handler.AddRole)
		apiRouter.GET("/role", handler.AllRole)
		apiRouter.GET("/role/:name", handler.RoleDetail)
	}

	//平台镜像 apis
	{
		apiRouter.GET("/image", handler.GetPlatformImageList)
		apiRouter.POST("/image", handler.AddPlatformImage)
	}
}
