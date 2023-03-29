package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
)

type UploadedFile struct {
	FileName string `json:"fileName"`
	Size     int64  `json:"size"` // 单位字节
	URL      string `json:"url"`
}

// GetFileInfo 获取文件信息
func GetFileInfo(fileId int64) (UploadedFile, error) {
	var fileInfo UploadedFile
	resultInfo, err := dal.GetFileInfo(fileId)
	if err != nil {
		log.Printf("[GetFileInfo] 服务获取文件信息失败")
		return fileInfo, errors.New("服务获取文件信息失败")
	}
	fileInfo.FileName = resultInfo.Name
	fileInfo.URL = resultInfo.Url
	fileInfo.Size = resultInfo.Size
	return fileInfo, nil
}
