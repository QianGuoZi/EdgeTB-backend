package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"fmt"
	"log"
	"os"
)

type PublicDataset struct {
	Id          int64  `form:"id" json:"id"`                   //编号
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Type        string `form:"type" json:"type"`               //数据集类型
	Description string `form:"description" json:"description"` //数据集描述
}

type PublicCheck struct {
	Keyword  string `form:"keyword" json:"keyword"`   //关键词搜索
	Type     string `form:"type" json:"type"`         //类型筛选
	PageSize int    `form:"pageSize" json:"pageSize"` //页大小
	PageNo   int    `form:"pageNo" json:"pageNo"`     //页码
}

// UploadDatasetFile 传入文件名、路径，获取类型、大小，返回文件id
func UploadDatasetFile(filePath, fileTypeStr string) (string, string, int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[UploadDatasetFile] 服务获取上传文件信息失败")
		return "", "", 0, errors.New("服务获取上传文件信息失败")
	}
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	return filePath, fi.Name(), int(fi.Size()), nil
}

// AddDatasetByUpload 添加数据集信息（upload）
func AddDatasetByUpload(setName string, description string, setType string, source string,
	url string, fileName string, setSize int, userName string) (int, error) {
	dataSet := dal.Dataset{}
	dataSet.DatasetName = setName
	dataSet.Description = description
	dataSet.Type = setType
	dataSet.State = 0 //0为公开数据集，1为私有状态
	dataSet.Source = source
	dataSet.Url = url
	dataSet.FileName = fileName
	dataSet.Size = int64(setSize)
	//通过username获取id
	userId, err := dal.GetUserId(userName)
	if err != nil {
		log.Printf("[AddDatasetByUpload] 服务获取用户id失败")
		return 0, errors.New("服务获取用户id失败")
	}
	dataSet.UserId = userId
	dataSetId, err := dal.AddDataset(dataSet)
	if err != nil {
		log.Printf("[AddDatasetByUpload] 服务通过上传添加数据集信息失败")
		return 0, errors.New("服务通过上传添加数据集信息失败")
	}
	return int(dataSetId), nil
}

// AllPublic 获取所有数据集List
func AllPublic(publicCheck PublicCheck) ([]PublicDataset, int, error) {
	offset := publicCheck.PageNo * publicCheck.PageSize
	datasets, err := dal.GetAllPublic(publicCheck.PageSize, offset)
	datasetTotal, err1 := dal.GetDatasetCount(0)
	var listLen int
	if datasetTotal < publicCheck.PageSize {
		listLen = datasetTotal
	} else {
		listLen = publicCheck.PageSize
	}
	publicDataset := make([]PublicDataset, listLen)
	if err != nil || err1 != nil {
		log.Printf("[AllPublicDatasets] 服务获取所有公共数据集列表失败")
		return publicDataset, 0, errors.New("服务获取所有公共数据集列表失败")
	}
	for i := 0; i < listLen; i++ {
		publicDataset[i].Id = datasets[i].Id
		publicDataset[i].SetName = datasets[i].DatasetName
		publicDataset[i].Type = datasets[i].Type
		publicDataset[i].Description = datasets[i].Description
	}
	log.Printf("[AllPublicDatasets] 服务获取所有公共数据集列表成功，内容为：%+v", publicDataset)

	return publicDataset, datasetTotal, nil
}
