package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
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

// UploadManager 上传manager脚本
func UploadManager(c *gin.Context) {
	//上传文件
	//上传到本地后添加新文件信息
	//获取用户名称
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}

	//处理文件上传
	//目录
	var saveDir = ""
	//名称
	var saveName = ""
	//完整路径
	var savePath = ""
	//获取文件
	file, errFile := c.FormFile("file")
	//处理获取文件错误
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择文件",
		})
		return
	}
	//获取上传文件的文件名和后缀(类型)
	uploadFileNameWithSuffix := path.Base(file.Filename) //得到完整文件名
	uploadFileType := path.Ext(uploadFileNameWithSuffix) //得到.txt等
	log.Printf("[UploadManager] fileName=%+v", file.Filename)
	log.Printf("[UploadManager] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法，应为py
	if uploadFileType != ".py" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件类型有误",
		})
		return
	}
	//文件保存目录
	saveDir = "./ProjectFile/Manager"
	//保存的文件名称
	dayStr := time.Now().String()
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadManager] savePath=%+v", savePath)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件上传失败",
		})
		return
	}
	//没有错误的情况下
	//返回文件path,size和fileName
	filePath, fileName, size, err1 := service.UploadManager(savePath)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息获取失败",
		})
		return
	}
	returnData := FileReturn{filePath, fileName, size}
	log.Printf("[UploadManager] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data":    returnData,
	})
	return
}

// UploadStructure 上传DML结构配置脚本
func UploadStructure(c *gin.Context) {
	//上传文件
	//上传到本地后添加新文件信息
	//获取用户名称
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}

	//处理文件上传
	//目录
	var saveDir = ""
	//名称
	var saveName = ""
	//完整路径
	var savePath = ""
	//获取文件
	file, errFile := c.FormFile("file")
	//处理获取文件错误
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择文件",
		})
		return
	}
	//获取上传文件的文件名和后缀(类型)
	uploadFileNameWithSuffix := path.Base(file.Filename) //得到完整文件名
	uploadFileType := path.Ext(uploadFileNameWithSuffix) //得到.txt等
	log.Printf("[UploadStructure] fileName=%+v", file.Filename)
	log.Printf("[UploadStructure] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法
	checkResult := CheckZipFile(uploadFileType)
	log.Printf("[UploadStructure] checkResult=%+v", checkResult)
	if checkResult != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件类型有误",
		})
		return
	}
	//文件保存目录
	saveDir = "./ProjectFile/Structure"
	//保存的文件名称
	dayStr := time.Now().String()
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadStructure] savePath=%+v", savePath)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件上传失败",
		})
		return
	}
	//没有错误的情况下
	//返回文件path,size和fileName
	filePath, fileName, size, err1 := service.UploadStructure(savePath)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息获取失败",
		})
		return
	}
	returnData := FileReturn{filePath, fileName, size}
	log.Printf("[UploadStructure] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data":    returnData,
	})
	return
}

// UploadDatasetSplitter 上传数据集切分脚本
func UploadDatasetSplitter(c *gin.Context) {
	//上传文件
	//上传到本地后添加新文件信息
	//获取用户名称
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}

	//处理文件上传
	//目录
	var saveDir = ""
	//名称
	var saveName = ""
	//完整路径
	var savePath = ""
	//获取文件
	file, errFile := c.FormFile("file")
	//处理获取文件错误
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择文件",
		})
		return
	}
	//获取上传文件的文件名和后缀(类型)
	uploadFileNameWithSuffix := path.Base(file.Filename) //得到完整文件名
	uploadFileType := path.Ext(uploadFileNameWithSuffix) //得到.txt等
	log.Printf("[UploadDatasetSplitter] fileName=%+v", file.Filename)
	log.Printf("[UploadDatasetSplitter] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法
	checkResult := CheckZipFile(uploadFileType)
	log.Printf("[UploadDatasetSplitter] checkResult=%+v", checkResult)
	if checkResult != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件类型有误",
		})
		return
	}
	//文件保存目录
	saveDir = "./ProjectFile/DatasetSplitter"
	//保存的文件名称
	dayStr := time.Now().String()
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadDatasetSplitter] savePath=%+v", savePath)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件上传失败",
		})
		return
	}
	//没有错误的情况下
	//返回文件path,size和fileName
	filePath, fileName, size, err1 := service.UploadDatasetSplitter(savePath)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息获取失败",
		})
		return
	}
	returnData := FileReturn{filePath, fileName, size}
	log.Printf("[UploadDatasetSplitter] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data":    returnData,
	})
	return
}

// AddProjectConfig 添加项目配置
func AddProjectConfig(c *gin.Context) {
	//通过用户id添加项目配置
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
	var configRequest service.ConfigRequest
	err = c.ShouldBind(&configRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "配置数据格式有误",
		})
		return
	}
	err = service.AddProjectConfig(username, projectName, configRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目配置添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目配置添加成功",
	})
	return
}

func GetProjectConfigs(c *gin.Context) {
	//获取项目配置
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
	configs, err := service.GetProjectConfigList(username, projectName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目配置获取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目配置获取成功",
		"data":    configs,
	})
}

// StartProject 运行项目
func StartProject(c *gin.Context) {
	//通过用户id添加项目配置
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
	log.Printf("[StartProject] projectName=%+v", projectName)
	err = service.StartProject(username, projectName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目运行失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目运行成功",
	})
	return
}

// FinishProject 终止项目
func FinishProject(c *gin.Context) {
	//通过用户id添加项目配置
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
	log.Printf("[StartProject] projectName=%+v", projectName)
	err = service.FinishProject(username, projectName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目终止失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目终止成功",
	})
	return
}
