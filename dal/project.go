package dal

import (
	"errors"
	"log"
)

// GetProjectId 获取项目id
func GetProjectId(projectName string) (int64, error) {
	project := Project{}
	DB.Model(&Project{}).Where("project_name = ?", projectName).First(&project)
	if project.Id == 0 {
		log.Printf("[GetProjectId] 无法找到该项目")
		return 0, errors.New("无法找到该项目")
	}
	log.Printf("[GetProjectId] projectId=%+v", project.Id)
	return project.Id, nil
}

// AddProject 添加项目信息
func AddProject(project Project) (int64, error) {
	result := DB.Model(&Project{}).Create(&project)
	if result.Error != nil {
		log.Printf("[AddProject] 数据库创建项目失败")
		return 0, result.Error
	}
	return project.Id, nil
}

// GetAllProject 获取项目列表
func GetAllProject(userId int64) ([]Project, error) {
	var projects []Project
	result := DB.Model(&Project{}).Where("user_id = ?", userId).Find(&projects)
	if result.Error != nil {
		log.Printf("[GetAllProject] 数据库获取项目列表信息失败")
		return projects, result.Error
	}
	listLen := len(projects)
	log.Printf("[GetAllProject] 数据库获取项目列表信息成功，长度为：%+v", listLen)
	return projects, nil
}

// GetProjectInfo 获取项目详细信息
func GetProjectInfo(userId int64, projectName string) (Project, error) {
	var projectInfo Project
	result := DB.Model(&Project{}).Where("user_id = ? && project_name = ?", userId, projectName).First(&projectInfo)
	if result.Error != nil {
		log.Printf("[GetProjectInfo] 数据库获取项目信息失败")
		return projectInfo, result.Error
	}
	log.Printf("[GetProjectInfo] 数据库获取项目信息成功")
	return projectInfo, nil
}
