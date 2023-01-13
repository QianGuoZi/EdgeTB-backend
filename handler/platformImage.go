package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

//GetPlatformImageList 获取平台镜像列表
func GetPlatformImageList(c *gin.Context) {
	returnData, err := service.GetPlatformImage()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "查询平台镜像列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取平台镜像列表成功",
		"data":    returnData,
	})
	return
}

// AddPlatformImage 添加平台镜像
func AddPlatformImage(c *gin.Context) {
	var platformImage service.PlatformImage
	err := c.ShouldBind(&platformImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}
	err = service.AddPlatformImage(platformImage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "平台镜像添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "平台镜像添加成功",
	})
	return
}
