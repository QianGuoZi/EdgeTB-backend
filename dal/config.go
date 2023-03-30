package dal

import (
	"errors"
	"log"
)

// CheckProjectConfig 检查项目是否拥有config
func CheckProjectConfig(projectId int64) error {
	config := Config{}
	DB.Model(&Config{}).Where("project_id = ?", projectId).First(&config)
	if config.Id == 0 {
		return errors.New("该项目下无法找到配置文件")
	}
	log.Printf("[GetRoleId] roleId=%+v", config.Id)
	return nil
}

// AddConfig 创建配置
func AddConfig(config Config) (int64, error) {
	result := DB.Model(&Config{}).Create(&config)
	if result.Error != nil {
		log.Printf("[AddConfig] 数据库创建配置失败")
		return 0, result.Error
	}
	return config.Id, nil
}
