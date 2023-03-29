package dal

import "log"

// AddFile 添加文件信息
func AddFile(fileInfo File) (int64, error) {
	result := DB.Model(&File{}).Create(&fileInfo)
	if result.Error != nil {
		log.Printf("[AddFileInfo] 数据库创建文件信息失败")
		return 0, result.Error
	}
	return fileInfo.Id, nil
}

// GetFileInfo 获取文件信息详情
func GetFileInfo(fileId int64) (File, error) {
	var fileInfo File
	result := DB.Model(&File{}).Where("id = ?", fileId).First(&fileInfo)
	if result.Error != nil {
		log.Printf("[GetFileInfo] 数据库获取文件信息详情失败")
		return fileInfo, result.Error
	}
	log.Printf("[GetFileInfo] 数据库获取文件信息详情成功")
	return fileInfo, nil
}
