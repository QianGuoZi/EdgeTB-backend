package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ProjectAddRequest struct {
	Name string `json:"name"` // 项目名称，唯一
}

// AddProject 添加项目
func AddProject(c *gin.Context) {
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
	var projectRequest ProjectAddRequest
	err = c.ShouldBind(&projectRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目数据格式有误",
		})
		return
	}
	err = service.AddProject(username, projectRequest.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目添加成功",
	})
	return
}

// AllProject 获取项目列表
func AllProject(c *gin.Context) {
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
	returnData, err1 := service.GetAllProject(username)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目列表获取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目列表获取成功",
		"data":    returnData,
	})
	return
}

// ProjectDetail 获取项目详情
func ProjectDetail(c *gin.Context) {
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
	projectName := c.Param("name")
	log.Printf("[GetProjectDetail] projectName=%+v", projectName)
	returnData, err1 := service.GetProjectDetail(username, projectName)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目详情获取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目详情获取成功",
		"data":    returnData,
	})
	return
}

// AddProjectInfo 添加项目内容
func AddProjectInfo(c *gin.Context) {
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
	projectName := c.Param("name")
	log.Printf("[GetProjectDetail] projectName=%+v", projectName)
	var projectRequest service.ProjectInfoRequest
	err = c.ShouldBind(&projectRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目数据格式有误",
		})
		return
	}
	err = service.AddProjectDetail(username, projectName, projectRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目内容添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目内容添加成功",
	})
	return
}
