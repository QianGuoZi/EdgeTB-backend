package dal

import (
	"errors"
	"log"
)

// AddDataset 添加新数据集记录
func AddDataset(dataset Dataset) (int64, error) {
	result := DB.Model(&Dataset{}).Create(&dataset)
	if result.Error != nil {
		log.Printf("[AddDataset] 数据库创建新数据集记录失败")
		return 0, result.Error
	}
	return dataset.Id, nil
}

// GetAllPublic 获取所有公共数据集列表
func GetAllPublic(pageSize int, offset int, keyword string, fileType string) ([]Dataset, int, error) {
	var datasets []Dataset
	//筛选出从第offset+1开始的pageSize条数据
	Db := DB
	if keyword != "" {
		Db = Db.Where("dataset_name LIKE ?", "%"+keyword+"%")
	}
	if fileType != "" {
		Db = Db.Where("type = ?", fileType)
	}
	result := Db.Model(&Dataset{}).Select("id", "dataset_name", "type", "description").
		Where("state = 0").Limit(pageSize).Offset(offset).Find(&datasets)
	listLen := len(datasets)
	if result.Error != nil {
		log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表失败")
		return datasets, listLen, errors.New("数据库获取所有公共数据集列表失败")
	}
	log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表成功，长度为：%+v", listLen)
	return datasets, listLen, nil
}

// GetDatasetCount 获取某状态数据集数量
func GetDatasetCount(state int, keyword string, fileType string) (int, error) {
	Db := DB
	if keyword != "" {
		Db = Db.Where("dataset_name LIKE ?", "%"+keyword+"%")
	}
	if fileType != "" {
		Db = Db.Where("type = ?", fileType)
	}
	var count int64
	result := Db.Model(&Dataset{}).Where("state = ?", state).Count(&count)
	if result.Error != nil {
		log.Printf("[GetDatasetCount] 数据库查询%+v状态数据集数量失败", state)
		return 0, result.Error
	}
	log.Printf("[GetDatasetCount] 数据库查询%+v状态数据集数量为%+v", state, count)
	return int(count), result.Error
}

// GetDatasetDetail 获取数据集详细内容
func GetDatasetDetail(datasetId int) (Dataset, error) {
	dataset := Dataset{}
	result := DB.Model(&Dataset{}).Where("id = ?", datasetId).First(&dataset)
	if result.Error != nil {
		log.Printf("[GetPublicDetail] 数据库查询数据集内容失败")
		return dataset, errors.New("数据库查询数据集内容失败")
	}
	log.Printf("[GetPublicDetail] 数据库查询数据集内容为%+v", dataset)
	return dataset, nil
}

// GetAllPrivate 获取用户自定义数据集列表
func GetAllPrivate(userId int) ([]Dataset, int, error) {
	var datasets []Dataset
	result := DB.Model(&Dataset{}).Where("user_id = ?", userId).Find(&datasets)
	listLen := len(datasets)
	if result.Error != nil {
		log.Printf("[GetAllPrivate] 数据库获取user_id:%+v用户自定义数据集列表失败", userId)
		return datasets, listLen, errors.New("数据库获取用户自定义数据集列表失败")
	}
	log.Printf("[GetAllPrivate] 数据库获取用户自定义数据集列表成功，长度为：%+v", listLen)
	return datasets, listLen, nil
}

// CheckUserDataset 检查userId和datasetId是否匹配
func CheckUserDataset(userId int, datasetId int) error {
	dataset := Dataset{}
	result := DB.Model(&Dataset{}).Where("id = ? && user_id = ?", datasetId, userId).First(&dataset)
	if result.Error != nil {
		log.Printf("[UpdateDataset] 数据库中用户与数据集信息不匹配")
		return result.Error
	}
	return nil
}

// UpdateDataset 更新数据集信息
func UpdateDataset(datasetId int, setName string, description string) error {
	result := DB.Model(&Dataset{}).Where("id = ?", datasetId).
		Updates(Dataset{DatasetName: setName, Description: description})
	if result.Error != nil {
		log.Printf("[UpdateDataset] 数据库更新数据集信息失败")
		return result.Error
	}
	return nil
}

// DeleteDataset 删除数据集信息
func DeleteDataset(datasetId int) error {
	result := DB.Delete(&Dataset{}, datasetId)
	if result.Error != nil {
		log.Printf("[DeleteDataset] 数据库删除数据集信息失败")
		return result.Error
	}
	return nil
}
