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

type PublicDetails struct {
	Id          int64  `form:"id" json:"id"`                   //编号
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Type        string `form:"type" json:"type"`               //数据集类型
	Description string `form:"description" json:"description"` //数据集描述
	Size        int    `form:"size" json:"size"`               //文件大小
}

type PrivateDataset struct {
	Id          int64  `form:"id" json:"id"`                   //编号
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Type        string `form:"type" json:"type"`               //数据集类型
	Description string `form:"description" json:"description"` //数据集描述
	CreatedAt   string `form:"createdAt" json:"createdAt"`     //创建时间
	State       int8   `form:"state" json:"state"`             //数据集是否公开状态
}

type PrivateDetails struct {
	Id          int64  `form:"id" json:"id"`                   //编号
	SetName     string `form:"setName" json:"setName"`         //数据集名称
	Type        string `form:"type" json:"type"`               //数据集类型
	Description string `form:"description" json:"description"` //数据集描述
	Size        int    `form:"size" json:"size"`               //文件大小
	CreatedAt   string `form:"createdAt" json:"createdAt"`     //创建时间
}

// UploadDatasetFile 传入文件名、路径，获取类型、大小
func UploadDatasetFile(filePath string) (string, string, int, error) {
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
	url string, fileName string, setSize int, username string) (int, error) {
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
	userId, err := dal.GetUserId(username)
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

// AllPublicDatasets 获取所有数据集List
func AllPublicDatasets(publicCheck PublicCheck) ([]PublicDataset, int, error) {
	offset := publicCheck.PageNo * publicCheck.PageSize
	//获取≤页面大小数量的数据集List
	datasets, datasetListLen, err := dal.GetAllPublic(publicCheck.PageSize, offset, publicCheck.Keyword, publicCheck.Type)
	datasetTotal, err1 := dal.GetDatasetCount(0, publicCheck.Keyword, publicCheck.Type)
	var listLen int
	//判断返回的数据集数量是否小于页面大小
	if datasetListLen < publicCheck.PageSize {
		listLen = datasetListLen
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

// GetPublicDatasetDetails 获取公共数据集详情
func GetPublicDatasetDetails(datasetId int) (PublicDetails, error) {
	publicDetails := PublicDetails{}
	details, err := dal.GetDatasetDetail(datasetId)
	if err != nil {
		log.Printf("[GetPublicDatasetDetails] 服务获取公共数据集详情失败")
		return publicDetails, errors.New("服务获取公共数据集详情失败")
	}
	publicDetails.Id = details.Id
	publicDetails.SetName = details.DatasetName
	publicDetails.Description = details.Description
	publicDetails.Type = details.Type
	publicDetails.Size = int(details.Size)
	log.Printf("[GetPublicDatasetDetails] 服务获取公共数据集详情为%+v", publicDetails)
	return publicDetails, nil
}

// AllPrivateDatasets 获取用户所有数据集List
func AllPrivateDatasets(username string) ([]PrivateDataset, error) {
	privateDataset := make([]PrivateDataset, 1)
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AllPrivateDatasets] 服务获取用户id失败")
		return privateDataset, errors.New("服务获取用户id失败")
	}
	datasets, listLen, err1 := dal.GetAllPrivate(int(userId))
	if err1 != nil {
		log.Printf("[AllPrivateDatasets] 服务获取用户自定义数据集列表失败")
		return privateDataset, errors.New("服务获取用户自定义数据集列表失败")
	}
	privateDataset = make([]PrivateDataset, listLen)
	for i := 0; i < listLen; i++ {
		privateDataset[i].Id = datasets[i].Id
		privateDataset[i].SetName = datasets[i].DatasetName
		privateDataset[i].Type = datasets[i].Type
		privateDataset[i].Description = datasets[i].Description
		privateDataset[i].State = datasets[i].State
		timeStr := datasets[i].CreatedAt.String()
		privateDataset[i].CreatedAt = timeStr[0:10]
	}
	log.Printf("[AllPrivateDatasets] 服务获取用户自定义数据集列表成功，内容为：%+v", privateDataset)
	return privateDataset, nil
}

// GetPrivateDatasetDetails 获取自定义数据集详情
func GetPrivateDatasetDetails(username string, datasetId int) (PrivateDetails, error) {
	privateDetails := PrivateDetails{}
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetPrivateDatasetDetails] 服务获取用户id失败")
		return privateDetails, errors.New("服务获取用户id失败")
	}
	//检查用户与数据集是否匹配
	err = dal.CheckUserDataset(int(userId), datasetId)
	if err != nil {
		log.Printf("[GetPrivateDatasetDetails] 服务用户与数据集id不匹配")
		return privateDetails, errors.New("服务用户与数据集id不匹配")
	}

	details, err1 := dal.GetDatasetDetail(datasetId)
	if err1 != nil {
		log.Printf("[GetPrivateDatasetDetails] 服务获取自定义数据集详情失败")
		return privateDetails, errors.New("服务获取自定义数据集详情失败")
	}
	privateDetails.Id = details.Id
	privateDetails.SetName = details.DatasetName
	privateDetails.Description = details.Description
	privateDetails.Type = details.Type
	privateDetails.Size = int(details.Size)
	timeStr := details.CreatedAt.String()
	privateDetails.CreatedAt = timeStr[0:19]
	log.Printf("[GetPrivateDatasetDetails] 服务获取自定义数据集详情为%+v", privateDetails)
	return privateDetails, nil
}

// UpdateDataset 修改数据集信息
func UpdateDataset(username string, datasetId int, setName string, description string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[UpdateDataset] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	err = dal.CheckUserDataset(int(userId), datasetId)
	if err != nil {
		log.Printf("[UpdateDataset] 服务用户与数据集id不匹配")
		return errors.New("服务用户与数据集id不匹配")
	}
	err = dal.UpdateDataset(datasetId, setName, description)
	if err != nil {
		log.Printf("[UpdateDataset] 服务修改数据集信息失败")
		return errors.New("服务修改数据集信息失败")
	}
	return nil
}

// DeleteDataset 删除数据集信息
func DeleteDataset(username string, datasetId int) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[DeleteDataset] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	err = dal.CheckUserDataset(int(userId), datasetId)
	if err != nil {
		log.Printf("[DeleteDataset] 服务用户与数据集id不匹配")
		return errors.New("服务用户与数据集id不匹配")
	}
	err = dal.DeleteDataset(datasetId)
	if err != nil {
		log.Printf("[DeleteDataset] 服务删除数据集信息失败")
		return errors.New("服务删除数据集信息失败")
	}
	return nil
}
