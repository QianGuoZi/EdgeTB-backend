package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
)

func AddProject(username, projectName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddDatasetByUpload] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	log.Printf(string(userId))
	return nil
}
