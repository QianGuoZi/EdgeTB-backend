package dal

import (
	"errors"
	"log"
)

// AddDatasetFile 添加新数据集文件记录
func AddDatasetFile(datasetFile DatasetFile) (int64, error) {
	result := DB.Model(&DatasetFile{}).Create(&datasetFile)
	if result.Error != nil {
		log.Printf("[AddDatasetFile] 数据库创建新数据集文件记录失败")
		return 0, result.Error
	}
	return datasetFile.Id, nil
}

// DeleteDatasetFile 删除数据集文件记录
func DeleteDatasetFile(fileId int64) error {
	result := DB.Delete(&DatasetFile{}, fileId)
	if result.Error != nil {
		log.Printf("[DeleteDatasetFile] 数据库删除数据集文件记录失败")
		return result.Error
	}
	return nil
}

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
func GetAllPublic(pageSize, offset int) ([]Dataset, int, error) {
	var datasets []Dataset
	//筛选出从第offset+1开始的pageSize条数据
	DB.Model(&Dataset{}).Select("id", "set_name", "type", "description").
		Where("state = 0").Limit(pageSize).Offset(offset).Find(&datasets)
	listLen := len(datasets)
	if listLen != 0 {
		log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表成功，长度为：%+v", listLen)
		return datasets, listLen, nil
	}
	log.Printf("[GetAllPublic] 数据库获取所有公共数据集列表失败")
	return datasets, listLen, errors.New("数据库获取所有公共数据集列表失败")
}

func GetAllPrivate() {

}

func GetPublicInfo() {

}
