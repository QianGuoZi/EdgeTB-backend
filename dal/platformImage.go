package dal

import (
	"errors"
	"log"
)

// 用于管理平台提供的镜像

// AddPlatformImage 添加新平台镜像
func AddPlatformImage(platformImage PlatformImage) (int64, error) {
	result := DB.Model(&PlatformImage{}).Create(&platformImage)
	if result.Error != nil {
		log.Printf("[AddPlatformImage] 数据库创建新平台镜像失败")
		return 0, result.Error
	}
	return platformImage.Id, nil
}

// GetAllPlatformImage 获取平台镜像列表
func GetAllPlatformImage() ([]PlatformImage, int, error) {
	var platformImages []PlatformImage
	result := DB.Model(&PlatformImage{}).Select("image_name", "description").Find(&platformImages)
	listLen := len(platformImages)
	if result.Error != nil {
		log.Printf("[GetAllPlatformImage] 数据库获取平台镜像列表失败")
		return platformImages, listLen, errors.New("数据库获取平台镜像列表失败")
	}
	log.Printf("[GetAllPrivate] 数据库获取平台镜像列表成功，长度为：%+v", listLen)
	return platformImages, listLen, nil
}
