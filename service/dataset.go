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

// UploadDatasetFile 传入文件名、路径，获取类型、大小，返回文件id
func UploadDatasetFile(filePath, fileTypeStr string) (int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, errors.New("")
	}
	var datasetFile dal.DatasetFile
	datasetFile.FileName = fi.Name()
	datasetFile.Size = fi.Size()
	datasetFile.Path = filePath
	datasetFile.DatasetId = 0
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	fmt.Println("fileTypeStr:", fileTypeStr)
	log.Printf("[UploadDatasetFile] 新的数据集文件信息为%+v", datasetFile)
	fileId, err := dal.AddDatasetFile(datasetFile)
	if err != nil {
		log.Printf("[UploadDatasetFile] 服务添加数据集文件信息失败")
		return 0, errors.New("服务添加数据集文件信息失败")
	}
	return int(fileId), nil
}

func AddDatasetFile() {

}

func AddDataset(setName string, setType string, setSize int, description string, state int, userId int) (int, error) {
	return 0, nil
}

//func AllPublic(publicCheck handler.PublicCheck) ([]PublicDataset, int, error) {
//	offset := publicCheck.PageNo * publicCheck.PageSize
//	datasets, listLen, err := dal.GetAllPublic(publicCheck.PageSize, offset)
//	publicDataset := make([]PublicDataset, listLen)
//	for i := 0; i < listLen; i++ {
//		publicDataset[i].Id = datasets[i].Id
//		publicDataset[i].SetName = datasets[i].SetName
//		publicDataset[i].Type = datasets[i].Type
//		publicDataset[i].Description = datasets[i].Description
//	}
//	log.Printf("[AllPublicDatasets] 服务获取所有公共数据集列表成功，内容为：%+v", publicDataset)
//	if err != nil {
//		log.Printf("[AllPublicDatasets] 服务获取所有公共数据集列表失败")
//		return publicDataset, 0, errors.New("服务获取所有公共数据集列表失败")
//	}
//	return publicDataset, listLen, nil
//}
