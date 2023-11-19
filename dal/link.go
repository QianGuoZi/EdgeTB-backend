package dal

import "log"

func AddLinks(links []Link) (int64, error) {
	result := DB.Model(&Link{}).Create(&links)
	if result.Error != nil {
		log.Printf("[AddLinks] 数据库创建link失败")
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func GetAllLinks(configId int64) ([]Link, error) {
	var links []Link
	result := DB.Model(&Link{}).Where("config_id = ?", configId).Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
