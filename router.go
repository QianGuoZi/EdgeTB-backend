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
		apiRouter.POST("/role/upload/code", handler.UploadRoleCode)
		apiRouter.PUT("/role/:name", handler.UpdateRole)
		apiRouter.DELETE("/role/:name", handler.DeleteRole)
	}

	//平台镜像 apis
	{
		apiRouter.GET("/image", handler.GetPlatformImageList)
		apiRouter.POST("/image", handler.AddPlatformImage)
	}

	//项目 apis
	{
		apiRouter.POST("/project", handler.AddProject)
		apiRouter.GET("/project", handler.AllProject)
		apiRouter.GET("/project/:name", handler.ProjectDetail)
		apiRouter.POST("/project/:name", handler.AddProjectInfo)
		apiRouter.POST("/project/:name/manager", handler.UploadManager)
		apiRouter.POST("/project/:name/structure_conf", handler.UploadStructure)
		apiRouter.POST("/project/:name/dataset_conf", handler.UploadDatasetSplitter)
		apiRouter.POST("/project/:name/config", handler.AddProjectConfig)
		apiRouter.GET("/project/:name/config", handler.GetProjectConfigs)
		apiRouter.GET("/project/:name/start", handler.StartProject)
		apiRouter.GET("/project/:name/finish", handler.FinishProject)
	}

	//日志 apis
	{
		apiRouter.POST("/log", handler.AddLog)
		apiRouter.GET("/log", handler.AllLog)
	}

	{
		apiRouter.POST("/task", handler.AddTask)
		apiRouter.GET("/task", handler.GetAllTasks)
		apiRouter.POST("/task/:id/start", handler.StartTask)
	}
}
