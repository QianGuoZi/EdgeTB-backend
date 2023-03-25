package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Project struct {
	Name string `json:"name"` // 项目名称，唯一
}

// AddProject 添加平台镜像
func AddProject(c *gin.Context) {
	//通过用户id获取自定义数据集列表
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	var projectRequest Project
	err = c.ShouldBind(&projectRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "数据格式有误",
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
