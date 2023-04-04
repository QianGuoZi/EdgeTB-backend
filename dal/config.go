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

// GetConfigInfo 获取配置信息
func GetConfigInfo(projectId int64) (Config, error) {
	project, err := GetProjectInfoById(projectId)
	if err != nil {
		return Config{}, err
	}
	var config Config
	// result := DB.Model(&Config{}).Where("project_id = ?", projectId).First(&config)
	result := DB.Model(&Config{}).Where("id = ?", project.CurrentConfigId).First(&config)
	if result.Error != nil {
		return config, result.Error
	}
	return config, nil
}

func GetAllConfig(projectId int64) ([]Config, error) {
	var configs []Config
	result := DB.Model(&Config{}).Where("project_id = ?", projectId).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

func GetConfigNodes(configId int64) ([]Node, error) {
	var nodeGroups []Node
	result := DB.Model(&Node{}).Where("config_id = ?", configId).Find(&nodeGroups)
	if result.Error != nil {
		return nil, result.Error
	}
	return nodeGroups, nil
}

func GetConfigNodeCount(configId int64) (int64, error) {
	var count int64
	result := DB.Model(&Node{}).Where("config_id = ?", configId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
