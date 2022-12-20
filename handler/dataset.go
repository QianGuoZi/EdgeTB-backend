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

type PublicCheck struct {
	Keyword  string `form:"keyword" json:"keyword"`   //关键词搜索
	Type     string `form:"type" json:"type"`         //类型筛选
	PageSize int    `form:"pageSize" json:"pageSize"` //页大小
	PageNo   int    `form:"pageNo" json:"pageNo"`     //页码
}

type PublicListReturn struct {
	Total    int                     `form:"total" json:"total"`       //数据集总数
	PageNo   int                     `form:"pageNo" json:"pageNo"`     //页码
	PageSize int                     `form:"pageSize" json:"pageSize"` //页大小
	Dataset  []service.PublicDataset `form:"dataset" json:"dataset"`   //数据集
}

type NewDataset struct {
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Description string `form:"description" json:"description"` //数据集描述
	Source      string `form:"source" json:"source"`           //数据集来源
	FileId      string `form:"fileId" json:"fileId"`           //已上传的文件标识
}

type DatasetFileReturn struct {
	FileId string `form:"fileId" json:"fileId"`
}

// AllPublicDatasets 公共数据集列表
//func AllPublicDatasets(c *gin.Context) {
//	//返回所有公共数据集（切页）
//	var publicCheck PublicCheck
//	err := c.ShouldBind(&publicCheck)
//	log.Printf("[AllPublicDatasets] publicCheck=%+v", publicCheck)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"success": false,
//			"message": "查询条件为空",
//		})
//		return
//	}
//
//	datasets, total, err1 := service.AllPublic(publicCheck)
//	if err1 != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"success": false,
//			"message": "查询公共数据集失败",
//		})
//		return
//	}
//	returnData := PublicListReturn{total, publicCheck.PageNo, publicCheck.PageSize, datasets}
//
//	c.JSON(http.StatusOK, gin.H{
//		"success": true,
//		"message": "登陆成功",
//		"data":    returnData,
//	})
//	return
//}

// PublicDatasetsDetail 公共数据集详情
func PublicDatasetsDetail() {
	//通过数据集id返回内容
	//id := c.Param("id")
}

// AllPrivateDatasets 自定义数据集列表
func AllPrivateDatasets() {
	//通过用户id获取自定义数据集列表
}

// PrivateDatasetsDetail 自定义数据集详情
func PrivateDatasetsDetail() {
	//通过数据集id返回内容
}

// DeletePrivateDataset 删除自定义数据集
func DeletePrivateDataset() {
	//拿到数据集id，删除对应内容
}

// AlterDataset 修改数据集
func AlterDataset() {
	//通过数据集id修改内容
}

// UploadDataset 上传数据集
func UploadDataset(c *gin.Context) {
	//上传文件
	//上传到本地后添加新文件信息

	//测试字段
	//var datasetFileReturn DatasetFileReturn
	//err := c.ShouldBind(&datasetFileReturn)
	//log.Printf("[UploadDataset] datasetFileReturn=%+v", datasetFileReturn)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"success": false,
	//		"message": "datasetFileReturn为空",
	//	})
	//	return
	//}
	//获取用户名称
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请选择文件",
		})
		return
	}
	//获取上传文件的文件名和后缀(类型)
	uploadFileNameWithSuffix := path.Base(file.Filename) //得到完整文件名
	uploadFileType := path.Ext(uploadFileNameWithSuffix) //得到.txt等
	log.Printf("[UploadDataset] fileName=%+v", file.Filename)
	log.Printf("[UploadDataset] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法
	checkResult := CheckFile(uploadFileType)
	if checkResult != true {
		if errFile != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "文件类型有误",
			})
			return
		}
	}
	//文件保存目录
	saveDir = "./DatasetFile"
	//保存的文件名称
	dayStr := time.Now().String()
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadDataset] savePath=%+v", savePath)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传失败",
		})
		return
	}
	//没有错误的情况下
	//传入文件名、路径、大小等文件信息，返回文件id
	fileId, err1 := service.UploadDatasetFile(savePath, uploadFileType)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息上传失败",
		})
		return
	}
	returnData := DatasetFileReturn{strconv.Itoa(fileId)}
	log.Printf("[UploadDataset] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传成功",
		"data":    returnData,
	})
	return
}

// CheckFile 检查文件类型是否符合
func CheckFile(uploadFileType string) bool {
	fileTypeList := []string{".zip", ".rar", ".gz", ".tar.gz", "tgz", "bz2", "z", "tar"}
	for _, element := range fileTypeList {
		if uploadFileType == element {
			return true
		}
	}
	return false
}

// AddDataset 创建自定义数据集
func AddDataset(c *gin.Context) {
	//通过用户id，文件id和数据集信息创建新数据集
	var newDataset NewDataset
	err := c.ShouldBind(&newDataset)
	log.Printf("[AddDataset] newDataset=%+v", newDataset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "新增数据集数据格式有误",
		})
		return
	}
	if newDataset.Source == "upload" {
		//datasets, total, err1 := service.AddDataset("数据集", "")
	}
}
