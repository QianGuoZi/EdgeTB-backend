package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
)

type PlatformImage struct {
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

// GetPlatformImage 获取平台镜像列表
func GetPlatformImage() ([]PlatformImage, error) {
	platformImageList := make([]PlatformImage, 1)
	result, listLen, err := dal.GetAllPlatformImage()
	if err != nil {
		log.Printf("[GetPlatformImage] 服务获取平台镜像列表失败")
		return platformImageList, errors.New("服务获取平台镜像列表失败")
	}
	platformImageList = make([]PlatformImage, listLen)
	for i := 0; i < listLen; i++ {
		platformImageList[i].Name = result[i].ImageName
		platformImageList[i].Description = result[i].Description
	}
	log.Printf("[GetPlatformImage] 服务获取平台镜像列表成功，内容为：%+v", platformImageList)
	return platformImageList, nil
}

// AddPlatformImage 添加平台镜像
func AddPlatformImage(platformImage PlatformImage) error {
	dataImage := dal.PlatformImage{}
	dataImage.ImageName = platformImage.Name
	dataImage.Description = platformImage.Description
	id, err := dal.AddPlatformImage(dataImage)
	if err != nil {
		log.Printf("[AddPlatformImage] 服务创建平台镜像失败")
		return errors.New("服务创建平台镜像失败")
	}
	log.Printf("[AddPlatformImage] 服务创建平台镜像成功，id为%+d", id)
	return nil
}
