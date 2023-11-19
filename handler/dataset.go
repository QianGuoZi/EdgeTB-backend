package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path"
	"strconv"
)

type PublicListReturn struct {
	Total    int                     `form:"total" json:"total"`       //数据集总数
	PageNo   int                     `form:"pageNo" json:"pageNo"`     //页码
	PageSize int                     `form:"pageSize" json:"pageSize"` //页大小
	Dataset  []service.PublicDataset `form:"dataset" json:"dataset"`   //数据集
}

type NewDataset struct {
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Description string `form:"description" json:"description"` //数据集描述
	Type        string `form:"type" json:"type"`               //数据集类型
	Source      string `form:"source" json:"source"`           //数据集来源
	Url         string `form:"url" json:"url"`                 //文件url
	FileName    string `form:"fileName" json:"fileName"`       //文件名称
	Size        int    `form:"size" json:"size"`               //文件大小
}

type DatasetFileReturn struct {
	Url      string `form:"url" json:"url"`           //文件url
	FileName string `form:"fileName" json:"fileName"` //文件名称
	Size     int    `form:"size" json:"size"`         //文件大小
}

type AddDatasetReturn struct {
	FileId int `form:"id" json:"id"` //数据集id
}

type DatasetDataUpdate struct {
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Description string `form:"description" json:"description"` //数据集描述
}

//AllPublicDatasets 公共数据集列表
func AllPublicDatasets(c *gin.Context) {
	//返回所有公共数据集（切页）
	var publicCheck service.PublicCheck
	err := c.ShouldBind(&publicCheck)
	if err != nil || publicCheck.PageSize == 0 {
		publicCheck.PageSize = 20
	}
	log.Printf("[AllPublicDatasets] publicCheck=%+v", publicCheck)
	datasets, total, err1 := service.AllPublicDatasets(publicCheck)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "查询公共数据集失败",
		})
		return
	}
	returnData := PublicListReturn{total, publicCheck.PageNo, publicCheck.PageSize, datasets}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取公共数据集列表成功",
		"data":    returnData,
	})
	return
}

// PublicDatasetsDetail 公共数据集详情
func PublicDatasetsDetail(c *gin.Context) {
	//通过数据集id返回内容
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "id有误，公共数据集详情获取失败",
		})
		return
	}
	log.Printf("[PublicDatasetsDetail] id=%+v", id)
	details, err1 := service.GetPublicDatasetDetails(id)
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "公共数据集详情获取失败，请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "公共数据集详情获取成功",
		"data":    details,
	})
	return
}

// AllPrivateDatasets 自定义数据集列表
func AllPrivateDatasets(c *gin.Context) {
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
	//返回datasetList
	datasets, err1 := service.AllPrivateDatasets(username)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "查询自定义数据集失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取自定义数据集列表成功",
		"data":    datasets,
	})
}

// PrivateDatasetsDetail 自定义数据集详情
func PrivateDatasetsDetail(c *gin.Context) {
	//通过数据集id返回内容
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	id, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "id有误，自定义数据集详情获取失败",
		})
		return
	}
	log.Printf("[PrivateDatasetsDetail] id=%+v", id)
	details, err2 := service.GetPrivateDatasetDetails(username, id)
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "自定义数据集详情获取失败，请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "自定义数据集详情获取成功",
		"data":    details,
	})
	return
}

// DeletePrivateDataset 删除自定义数据集
func DeletePrivateDataset(c *gin.Context) {
	//获取username
	username, err1 := service.GetUsername(c)
	if err1 != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err1)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	//拿到数据集id，删除对应内容
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "id有误，自定义数据集详情获取失败",
		})
		return
	}
	result := service.DeleteDataset(username, id)
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "数据集删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集删除成功",
	})
	return
}

// UpdateDataset 修改数据集
func UpdateDataset(c *gin.Context) {
	//通过数据集id修改内容
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "id有误，自定义数据集详情获取失败",
		})
		return
	}
	//获取username
	username, err1 := service.GetUsername(c)
	if err1 != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err1)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	var datasetDataUpdate DatasetDataUpdate
	err = c.ShouldBind(&datasetDataUpdate)
	log.Printf("[UpdateDataset] DatasetDataUpdate=%+v", datasetDataUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}
	result := service.UpdateDataset(username, id, datasetDataUpdate.SetName, datasetDataUpdate.Description)
	if result != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据集信息修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集信息修改成功",
	})
	return
}

// UploadDataset 上传数据集
func UploadDataset(c *gin.Context) {
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
	log.Printf("[UploadDataset] fileName=%+v", file.Filename)
	log.Printf("[UploadDataset] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法
	checkResult := CheckZipFile(uploadFileType)
	log.Printf("[UploadDataset] checkResult=%+v", checkResult)
	if checkResult != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件类型有误",
		})
		return
	}
	//文件保存目录
	saveDir = "./DatasetFile"
	//保存的文件名称
	//dayStr := time.Now().String()
	//timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	//saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	log.Printf(username)
	saveName = file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadDataset] savePath=%+v", savePath)
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
	filePath, fileName, size, err1 := service.UploadDatasetFile(savePath)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息获取失败",
		})
		return
	}
	returnData := DatasetFileReturn{filePath, fileName, size}
	log.Printf("[UploadDataset] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data":    returnData,
	})
	return
}

// AddDataset 创建自定义数据集
func AddDataset(c *gin.Context) {
	//通过用户id，文件信息和数据集信息创建新数据集

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

	//添加数据集信息
	var newDataset NewDataset
	err1 := c.ShouldBind(&newDataset)
	log.Printf("[AddDataset] newDataset=%+v", newDataset)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "新增数据集数据格式有误",
		})
		return
	}
	if newDataset.Source == "upload" {
		datasetId, err2 := service.AddDatasetByUpload(newDataset.SetName, newDataset.Description, newDataset.Type, newDataset.Source,
			newDataset.Url, newDataset.FileName, newDataset.Size, username)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "新增数据集失败，请重试",
			})
			return
		}
		returnData := AddDatasetReturn{datasetId}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "新增数据集成功",
			"data":    returnData,
		})
		return
	}
}
