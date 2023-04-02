package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
)

type LogRequest struct {
	Content  string `json:"content"`  // 日志内容
	NodeName string `json:"nodeName"` // 节点名称
}

type LogList struct {
	NodeName  string `json:"nodeName"`  // 节点名称
	Content   string `json:"content"`   // 内容
	CreatedAt string `json:"createdAt"` // 日志创建时间
}

func AddLog(projectName string, logRequest LogRequest) error {
	projectId, err1 := dal.GetProjectIdByName(projectName)
	if err1 != nil {
		log.Printf("[AddLog] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	var logInfo dal.Log
	logInfo.ProjectId = projectId
	logInfo.NodeName = logRequest.NodeName
	logInfo.Content = logRequest.Content
	_, err2 := dal.AddLog(logInfo)
	if err2 != nil {
		log.Printf("[AddLog] 服务创建日志失败")
		return errors.New("服务创建日志失败")
	}
	return nil
}

func AllLog(username, projectName string) ([]LogList, error) {
	var logL []LogList
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AllLog] 服务获取用户id失败")
		return logL, errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[AllLog] 服务获取项目id失败")
		return logL, errors.New("服务获取项目id失败")
	}
	logList, err2 := dal.GetLogList(projectId)
	if err2 != nil {
		log.Printf("[AllLog] 服务获取日志列表失败")
		return logL, errors.New("服务获取日志列表失败")
	}
	listLen := len(logList)
	logInfoList := make([]LogList, listLen)
	for i := 0; i < listLen; i++ {
		logInfoList[i].NodeName = logList[i].NodeName
		logInfoList[i].Content = logList[i].Content
		timeStr := logList[i].CreatedAt.String()
		logInfoList[i].CreatedAt = timeStr[0:19]
	}
	return logInfoList, nil
}
