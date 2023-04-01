package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
	
	projectName, _ := c.GetQuery("project")
	log.Printf("[AddTask] projectName=%+v", projectName)

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

	projectName, _ := c.GetQuery("project")
	log.Printf("[AddTask] projectName=%+v", projectName)

	var newTask service.AddTaskRequest
	err = c.ShouldBind(&newTask)
	log.Print(newTask)
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

	if has, err := dal.CheckIfProjectHasRunningTask(task.ProjectId); err == nil {
		if has {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "项目已有任务在运行",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取项目任务失败",
		})
		return
	}

	// 奇技淫巧！！
	// 任务运行将会修改项目的dataset_id和dataset_splitter_file_id
	// 然后就交给项目自己的运行了
	project := dal.Project{}
	err = dal.DB.Model(&dal.Project{}).Where("id = ?", task.ProjectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取项目失败",
		})
		return
	}
	project.CurrentConfigId = task.Id
	project.DatasetId = task.DatasetId
	project.DatasetSplitterFileId = task.DatasetSplitterFileId
	dal.DB.Save(&project)

	err = service.StartProject(username, project.ProjectName)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "启动任务失败",
		})
		return
	}

	task.Status = dal.TaskStatusRunning
	t := time.Now()
	task.StartAt = &t
	dal.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "启动任务成功",
	})
}

func FinishTask(c *gin.Context) {
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

	project := dal.Project{}
	err = dal.DB.Model(&dal.Project{}).Where("id = ?", task.ProjectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取项目失败",
		})
		return
	}

	err = service.FinishProject(username, project.ProjectName)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "停止任务失败",
		})
		return
	}

	task.Status = dal.TaskStatusStopped
	dal.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "停止任务成功",
	})
}
