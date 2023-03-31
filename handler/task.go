package handler

import (
	"log"
	"net/http"
	"strconv"

	"EdgeTB-backend/dal"
	"EdgeTB-backend/service"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	//获取用户username
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	log.Printf("[GetUserInfo] success username=%+v", username)

	projectName := c.Query("project")

	tasks, err := service.GetAllTask(username, projectName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取任务列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取任务列表成功",
		"data":    tasks,
	})
}

func AddTask(c *gin.Context) {
	//获取用户username
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	log.Printf("[GetUserInfo] success username=%+v", username)

	projectName := c.Query("project")

	var newTask service.AddTaskRequest
	err = c.ShouldBind(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "新增任务数据格式有误",
		})
		return
	}

	_, err = service.AddTask(username, projectName, newTask)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "新增任务失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "新增任务成功",
	})
}

func StartTask(c *gin.Context) {
	//获取用户username
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	log.Printf("[GetUserInfo] success username=%+v", username)

	taskIdStr := c.Param("id")
	taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "任务id有误",
		})
		return
	}

	task, err := dal.GetTask(taskId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取任务失败",
		})
		return
	}

	// 奇技淫巧！！
	// 任务运行将会修改项目的dataset_id和dataset_splitter_file_id
	// 然后就交给项目自己的运行了
	err = dal.DB.Model(&dal.Project{}).Where("id = ?", task.ProjectId).Updates(map[string]interface{}{
		"dataset_id":               task.DatasetId,
		"dataset_splitter_file_id": task.DatasetSplitterFileId,
	}).Error

	// TODO 任务运行

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "启动任务失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "启动任务成功",
	})
}
