package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"fmt"
	"log"
	"os"
)

type ProjectListResponse struct {
	Name string `json:"name"` // 项目名称，唯一
}

type ProjectInfoResponse struct {
	Config     bool        `json:"config"`               // 节点配置
	Controller *Controller `json:"controller,omitempty"` // 控制器
	Dataset    *Dataset    `json:"dataset,omitempty"`    // 数据集
	Name       string      `json:"name"`                 // 项目名称
}

type ProjectInfoRequest struct {
	Controller Controller `json:"controller"` // 控制器
	Dataset    Dataset    `json:"dataset"`    // 数据集
}

// Controller 控制器
type Controller struct {
	Manager   UploadedFile `json:"manager"`   // Manager.py文件信息
	Structure UploadedFile `json:"structure"` // DML结构配置文件信息
}

// Dataset 数据集
type Dataset struct {
	ID       int          `json:"id"`       // 选择的数据集id
	Splitter UploadedFile `json:"splitter"` // 数据集切分脚本文件信息
}

type ConfigRequest struct {
	Nodes          []NodeInfo `json:"nodes"`          // 节点数组
	Topology       string     `json:"topology"`       // 链路类型
	BandwidthLower int        `json:"bandwidthLower"` // 带宽下界（包含），单位mbps
	BandwidthUpper int        `json:"bandwidthUpper"` // 带宽上界（包含），单位mbps
}

type NodeInfo struct {
	Name string `json:"name"` // 节点名称
	CPU  int    `json:"cpu"`  // CPU大小
	RAM  int    `json:"ram"`  // RAM大小，单位MB
	Role string `json:"role"` // 角色名称
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
		log.Printf("[GetAllProject] 服务获取用户id失败")
		return list, errors.New("服务获取用户id失败")
	}
	projectList, err1 := dal.GetAllProject(userId)
	if err1 != nil {
		log.Printf("[GetAllProject] 服务获取项目列表失败")
		return list, errors.New("服务获取项目列表失败")
	}
	listLen := len(projectList)
	responseList := make([]ProjectListResponse, listLen)
	for i := 0; i < listLen; i++ {
		responseList[i].Name = projectList[i].ProjectName
	}
	return responseList, nil
}

// GetProjectDetail 获取项目详情
func GetProjectDetail(username, projectName string) (ProjectInfoResponse, error) {
	var projectResponse ProjectInfoResponse
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetProjectDetail] 服务获取用户id失败")
		return projectResponse, errors.New("服务获取用户id失败")
	}
	projectInfo, err1 := dal.GetProjectInfo(userId, projectName)
	if err1 != nil {
		log.Printf("[GetProjectDetail] 服务获取项目信息失败")
		return projectResponse, errors.New("服务获取项目信息失败")
	}

	//projectName
	projectResponse.Name = projectInfo.ProjectName
	//config
	checkResult := dal.CheckProjectConfig(projectInfo.Id)
	if checkResult != nil {
		projectResponse.Config = false
	} else {
		projectResponse.Config = true
	}
	//controller
	projectResponse.Controller = new(Controller)
	//manager文件信息
	if projectInfo.ManagerFileId == 0 {
		projectResponse.Controller = nil
	} else {
		managerFileInfo, _ := GetFileInfo(projectInfo.ManagerFileId)
		projectResponse.Controller.Manager = managerFileInfo
	}
	//structure文件信息
	if projectInfo.StructureFileId == 0 {
		projectResponse.Controller = nil
	} else {
		structureFileInfo, _ := GetFileInfo(projectInfo.StructureFileId)
		projectResponse.Controller.Structure = structureFileInfo
	}
	//dataset
	projectResponse.Dataset = new(Dataset)
	//datasetId
	if projectInfo.DatasetId == 0 {
		projectResponse.Dataset = nil
	} else {
		projectResponse.Dataset.ID = int(projectInfo.DatasetId)
	}
	//datasetSplitter文件信息
	if projectInfo.DatasetSplitterFileId == 0 {
		projectResponse.Dataset = nil
	} else {
		datasetSplitterFileInfo, _ := GetFileInfo(projectInfo.DatasetSplitterFileId)
		projectResponse.Dataset.Splitter = datasetSplitterFileInfo
	}
	return projectResponse, nil
}

// AddProjectDetail 添加项目内容
func AddProjectDetail(username, projectName string, projectInfo ProjectInfoRequest) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddProjectDetail] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	projectId, err := dal.GetProjectId(projectName, userId)
	if err != nil {
		log.Printf("[AddProjectDetail] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	//controller 信息添加
	//Manager
	var managerFile dal.File
	managerFile.Name = projectInfo.Controller.Manager.FileName
	managerFile.Url = projectInfo.Controller.Manager.URL
	managerFile.Size = projectInfo.Controller.Manager.Size
	managerId, err1 := dal.AddFile(managerFile)
	if err1 != nil {
		log.Printf("[AddProjectDetail] 服务添加manager文件信息失败")
		return errors.New("服务添加manager文件信息失败")
	}
	//structure
	var structureFile dal.File
	structureFile.Name = projectInfo.Controller.Structure.FileName
	structureFile.Url = projectInfo.Controller.Structure.URL
	structureFile.Size = projectInfo.Controller.Structure.Size
	structureId, err2 := dal.AddFile(structureFile)
	if err2 != nil {
		log.Printf("[AddProjectDetail] 服务添加structure文件信息失败")
		return errors.New("服务添加structure文件信息失败")
	}
	//dataset 信息添加
	//datasetSplitter
	var datasetSplitterFile dal.File
	datasetSplitterFile.Name = projectInfo.Dataset.Splitter.FileName
	datasetSplitterFile.Url = projectInfo.Dataset.Splitter.URL
	datasetSplitterFile.Size = projectInfo.Dataset.Splitter.Size
	datasetSplitterId, err3 := dal.AddFile(datasetSplitterFile)
	if err3 != nil {
		log.Printf("[AddProjectDetail] 服务添加datasetSplitter文件信息失败")
		return errors.New("服务添加datasetSplitter文件信息失败")
	}
	//存储4个id
	err = dal.UpdateProjectInfo(projectId, managerId, structureId, int64(projectInfo.Dataset.ID), datasetSplitterId)
	if err != nil {
		log.Printf("[AddProjectDetail] 服务添加项目信息失败")
		return errors.New("服务添加项目信息失败")
	}
	return nil
}

// UploadManager 传入文件名、路径，获取类型、大小
func UploadManager(filePath string) (string, string, int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[UploadManager] 服务获取上传文件信息失败")
		return "", "", 0, errors.New("服务获取上传文件信息失败")
	}
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	return filePath, fi.Name(), int(fi.Size()), nil
}

// UploadStructure 传入文件名、路径，获取类型、大小
func UploadStructure(filePath string) (string, string, int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[UploadStructure] 服务获取上传文件信息失败")
		return "", "", 0, errors.New("服务获取上传文件信息失败")
	}
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	return filePath, fi.Name(), int(fi.Size()), nil
}

// UploadDatasetSplitter 传入文件名、路径，获取类型、大小
func UploadDatasetSplitter(filePath string) (string, string, int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[UploadDatasetSplitter] 服务获取上传文件信息失败")
		return "", "", 0, errors.New("服务获取上传文件信息失败")
	}
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	return filePath, fi.Name(), int(fi.Size()), nil
}

// AddProjectConfig 添加项目配置
func AddProjectConfig(username, projectName string, configInfo ConfigRequest) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddProjectConfig] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[AddProjectConfig] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	//创建config
	var config dal.Config
	config.LinkType = configInfo.Topology
	config.BandwidthUpper = int64(configInfo.BandwidthUpper)
	config.BandwidthLower = int64(configInfo.BandwidthLower)
	config.ProjectId = projectId
	configId, err2 := dal.AddConfig(config)
	if err2 != nil {
		log.Printf("[AddProjectConfig] 服务创建配置失败")
		return errors.New("服务创建配置失败")
	}
	//创建node
	nodeLen := len(configInfo.Nodes)
	var node dal.Node
	for i := 0; i < nodeLen; i++ {
		node.ConfigId = configId
		node.NodeName = configInfo.Nodes[i].Name
		node.CPU = int64(configInfo.Nodes[i].CPU)
		node.RAM = int64(configInfo.Nodes[i].RAM)
		node.RoleName = configInfo.Nodes[i].Role
		_, err3 := dal.AddNode(node)
		if err3 != nil {
			log.Printf("[AddProjectConfig] 服务创建节点配置失败")
			return errors.New("服务创建节点配置失败")
		}
	}
	return nil
}

// StartProject 运行项目
func StartProject(username, projectName string) error {
	log.Printf(username, projectName)
	return nil
}

// FinishProject 终止项目
func FinishProject(username, projectName string) error {
	log.Printf(username, projectName)
	return nil
}

type ConfigResponse struct {
	Id        int64  `json:"id"`
	Topology  string `json:"topology"`
	NodeCount int    `json:"nodeCount"`
	CreatedAt string `json:"createdAt"`
}

func GetProjectConfigs(username, projectName string) ([]ConfigResponse, error) {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetProjectConfigs] 服务获取用户id失败")
		return nil, errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[GetProjectConfigs] 服务获取项目id失败")
		return nil, errors.New("服务获取项目id失败")
	}
	//获取config
	configs, err2 := dal.GetAllConfig(projectId)
	if err2 != nil {
		log.Printf("[GetProjectConfigs] 服务获取配置失败")
		return nil, errors.New("服务获取配置失败")
	}
	//获取node
	var configResponses []ConfigResponse
	for _, config := range configs {
		var configResponse ConfigResponse
		configResponse.Id = config.Id
		configResponse.Topology = config.LinkType
		configResponse.CreatedAt = config.CreatedAt.Format("2006-01-02 15:04:05")
		nodeCount, err := dal.GetConfigNodeCount(config.Id)
		if err != nil {
			log.Printf("[GetProjectConfigs] 服务获取节点数量失败")
			return nil, errors.New("服务获取节点数量失败")
		}
		configResponse.NodeCount = int(nodeCount)
		configResponses = append(configResponses, configResponse)
	}
	return configResponses, nil
}
