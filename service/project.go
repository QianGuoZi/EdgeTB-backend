package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
)

type ProjectListResponse struct {
	Name string `json:"name"` // 项目名称，唯一
}

// AddProject 添加项目
func AddProject(username, projectName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddProject] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	project := dal.Project{}
	project.ProjectName = projectName
	project.UserId = userId
	id, err1 := dal.AddProject(project)
	if err1 != nil {
		log.Printf("[AddProject] 服务添加项目失败")
		return errors.New("服务添加项目失败")
	}
	log.Printf("[AddProject] 服务添加项目成功，id为：%d", id)
	return nil
}

// GetAllProject 获取项目列表
func GetAllProject(username string) ([]ProjectListResponse, error) {
	var list []ProjectListResponse
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddProject] 服务获取用户id失败")
		return list, errors.New("服务获取用户id失败")
	}
	projectList, err1 := dal.GetAllProject(userId)
	if err1 != nil {
		log.Printf("[GetAllRole] 服务获取用户角色列表失败")
		return list, errors.New("服务获取用户角色列表失败")
	}
	listLen := len(projectList)
	responseList := make([]ProjectListResponse, listLen)
	for i := 0; i < listLen; i++ {
		responseList[i].Name = projectList[i].ProjectName
	}
	return responseList, nil
}

// GetProjectDetail 获取项目详情
func GetProjectDetail(username, projectName string) (ProjectListResponse, error) {
	var project ProjectListResponse
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddProject] 服务获取用户id失败")
		return project, errors.New("服务获取用户id失败")
	}
	projectInfo, err1 := dal.GetProjectInfo(userId, projectName)
	if err1 != nil {
		log.Printf("[GetProjectDetail] 服务获取项目信息失败")
		return project, errors.New("服务获取项目信息失败")
	}
	project.Name = projectInfo.ProjectName
	return project, nil
}
