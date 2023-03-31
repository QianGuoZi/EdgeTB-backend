package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AllLog 获取日志
func AllLog(c *gin.Context) {
	//通过用户id获取项目列表
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	//获取项目名称
	projectName, _ := c.GetQuery("project")
	log.Printf("[AllLog] projectName=%+v", projectName)
	returnData, err1 := service.AllLog(username, projectName)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "日志获取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "日志获取成功",
		"data":    returnData,
	})
	return
}

// AddLog 添加日志
func AddLog(c *gin.Context) {
	//通过用户id创建项目
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	//获取项目名称
	projectName, _ := c.GetQuery("project")
	log.Printf("[AddLog] projectName=%+v", projectName)

	var logRequest service.LogRequest
	err = c.ShouldBind(&logRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "日志数据格式有误",
		})
		return
	}
	err = service.AddLog(username, projectName, logRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "日志添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "日志添加成功",
	})
	return
}
