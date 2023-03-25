package dal

import "log"

// AddProject 添加项目信息
func AddProject(project Project) (int64, error) {
	result := DB.Model(&Project{}).Create(&project)
	if result.Error != nil {
		log.Printf("[AddProject] 数据库创建项目失败")
		return 0, result.Error
	}
	return project.Id, nil
}
