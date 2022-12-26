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
func GetAllPublic(pageSize, offset int) ([]Dataset, error) {
	var datasets []Dataset
	//筛选出从第offset+1开始的pageSize条数据
	result := DB.Model(&Dataset{}).Select("id", "dataset_name", "type", "description").
		Where("state = 0").Limit(pageSize).Offset(offset).Find(&datasets)
	listLen := len(datasets)
	if result.Error != nil {
		log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表失败")
		return datasets, errors.New("数据库获取所有公共数据集列表失败")
	}
	log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表成功，长度为：%+v", listLen)
	return datasets, nil
}

// GetDatasetCount 获取某状态数据集数量
func GetDatasetCount(state int) (int, error) {
	var count int64
	result := DB.Model(&Dataset{}).Where("state = ?", state).Count(&count)
	if result.Error != nil {
		log.Printf("[GetDatasetCount] 数据库查询%+v状态数据集数量失败", state)
		return 0, result.Error
	}
	log.Printf("[GetDatasetCount] 数据库查询%+v状态数据集数量为%+v", state, count)
	return int(count), result.Error
}

func GetAllPrivate() {

}

func GetPublicInfo() {

}
