package dal

import "log"

// AddLog 创建节点配置
func AddLog(logInfo Log) (int64, error) {
	result := DB.Model(&Log{}).Create(&logInfo)
	if result.Error != nil {
		log.Printf("[AddLog] 数据库创建日志失败")
		return 0, result.Error
	}
	return logInfo.Id, nil
}

// GetLogList 获取日志详情列表
func GetLogList(projectId int64) ([]Log, error) {
	var logList []Log
	result := DB.Model(&Log{}).Where("project_id = ?", projectId).Find(&logList)
	if result.Error != nil {
		log.Printf("[GetLogInfo] 数据库获取日志信息失败")
		return logList, result.Error
	}
	log.Printf("[GetLogInfo] 数据库获取日志信息成功")
	return logList, nil
}
