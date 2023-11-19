package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConfigRequest struct {
	Nodes []NodeInfo                `json:"nodes"` // 节点数组
	Links map[string]map[string]int `json:"links"` // 链路，source->target->bandwidth(单位mbps)
}

type NodeInfo struct {
	Name string `json:"name"` // 节点名称
	CPU  int    `json:"cpu"`  // CPU大小
	RAM  int    `json:"ram"`  // RAM大小，单位MB
	Role string `json:"role"` // 角色名称
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
	//创建link
	links := make([]dal.Link, 0)
	for source, target := range configInfo.Links {
		for targetName, bandwidth := range target {
			vlink := dal.Link{
				ConfigId:       configId,
				SourceNodeName: source,
				TargetNodeName: targetName,
				Bandwidth:      bandwidth,
			}
			links = append(links, vlink)
		}
	}
	_, err = dal.AddLinks(links)
	if err != nil {
		log.Printf("[AddProjectConfig] 服务创建link失败")
		return errors.New("服务创建link失败")
	}
	return nil
}

type ConfigResponse struct {
	Id        int64                     `json:"id"`
	NodeCount int                       `json:"nodeCount"`
	Links     map[string]map[string]int `json:"links"`
	CreatedAt string                    `json:"createdAt"`
}

// GetProjectConfigList 获取项目的配置
func GetProjectConfigList(username, projectName string) ([]ConfigResponse, error) {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetProjectConfigList] 服务获取用户id失败")
		return nil, errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[GetProjectConfigList] 服务获取项目id失败")
		return nil, errors.New("服务获取项目id失败")
	}
	//获取config
	configs, err2 := dal.GetAllConfig(projectId)
	if err2 != nil {
		log.Printf("[GetProjectConfigList] 服务获取配置失败")
		return nil, errors.New("服务获取配置失败")
	}
	var configResponses []ConfigResponse
	for _, config := range configs {
		var configResponse ConfigResponse
		configResponse.Id = config.Id
		configResponse.CreatedAt = config.CreatedAt.Format("2006-01-02 15:04:05")
		//获取link
		links, err := dal.GetAllLinks(config.Id)
		if err != nil {
			log.Printf("[GetProjectConfigList] 服务获取link失败")
			return nil, errors.New("服务获取link失败")
		}
		configResponse.Links = ConvertLinkArrayToMap(links)
		//获取node
		nodeCount, err := dal.GetConfigNodeCount(config.Id)
		if err != nil {
			log.Printf("[GetProjectConfigList] 服务获取节点数量失败")
			return nil, errors.New("服务获取节点数量失败")
		}
		configResponse.NodeCount = int(nodeCount)
		configResponses = append(configResponses, configResponse)
	}
	return configResponses, nil
}

type ConfigYamlStruct struct {
	Links      ConfigYamlLink            `yaml:"links"`
	Deployment ConfigYamlDeployment      `yaml:"deployment"`
	Roles      map[string]ConfigYamlRole `yaml:"roles"`
}

type ConfigYamlLink map[string]map[string]string // source->target->bandwidthmbps

type ConfigYamlDeployment struct {
	Emulated map[string]ConfigYamlNode `yaml:"emulated"`
}

type ConfigYamlNode struct {
	Role string `yaml:"role"`
	Cpu  int    `yaml:"cpu"`
	Ram  int    `yaml:"ram"`
}

type ConfigYamlRole struct {
	Name    string `yaml:"name"`
	RootDir string `yaml:"root_dir"`
	Cmd     string `yaml:"cmd"`
}

// CreateConfigYaml 创建yaml文件
func CreateConfigYaml(username, projectName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[CreateConfigYaml] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[CreateConfigYaml] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	configInfo, err2 := dal.GetConfigInfo(projectId)
	if err2 != nil {
		log.Printf("[CreateConfigYaml] 服务获取项目配置信息失败")
		return errors.New("服务获取项目配置信息失败")
	}
	links, err3 := dal.GetAllLinks(configInfo.Id)
	if err3 != nil {
		log.Printf("[CreateConfigYaml] 服务获取项目配置link信息失败")
		return errors.New("服务获取项目配置link信息失败")
	}

	configYamlLink := make(ConfigYamlLink)
	for source, target := range ConvertLinkArrayToMap(links) {
		configYamlLink[source] = make(map[string]string)
		for targetName, bandwidth := range target {
			configYamlLink[source][targetName] = fmt.Sprintf("%dmbps", bandwidth)
		}
	}
	configYamlNode := make(map[string]ConfigYamlNode)
	emulatedNodeList, err3 := dal.GetConfigNodes(configInfo.Id)
	if err3 != nil {
		log.Printf("[CreateConfigYaml] 服务获取项目配置节点信息失败")
		return errors.New("服务获取项目配置节点信息失败")
	}
	for _, node := range emulatedNodeList {
		configYamlNode[node.NodeName] = ConfigYamlNode{
			Role: node.RoleName,
			Cpu:  int(node.CPU),
			Ram:  int(node.RAM),
		}
	}
	emulatedYaml := ConfigYamlDeployment{
		Emulated: configYamlNode,
	}

	roles, err := dal.GetAllRoleByProjectId(projectId)
	if err != nil {
		log.Printf("[CreateConfigYaml] 服务获取项目配置角色信息失败")
		return errors.New("服务获取项目配置角色信息失败")
	}
	roleYaml := make(map[string]ConfigYamlRole)
	for _, role := range roles {
		roleYaml[role.RoleName] = ConfigYamlRole{
			Name:    role.RoleName,
			RootDir: role.WorkDir,
			Cmd:     role.RunCommand,
		}
	}

	configYaml := ConfigYamlStruct{
		Links:      configYamlLink,
		Deployment: emulatedYaml,
		Roles:      roleYaml,
	}
	//fileName := username + "_" + projectName + "_" + strconv.Itoa(int(configInfo.Id)) + "_config.yaml"
	err = GenerateYaml("config.yaml", "EdgeTB", configYaml)
	if err != nil {
		log.Printf("[CreateConfigYaml] 服务生成yaml文件失败")
		return errors.New("服务生成yaml文件失败")
	}
	return nil
}

// GenerateYaml 生成yaml文件 yamlFile：文件路径 key：yaml 的key value: yaml 的value
func GenerateYaml(yamlFile, key string, value interface{}) error {
	filename := yamlFile
	var viperObj = viper.New()
	viperObj.SetConfigFile(filename)
	viperObj.SetConfigType("yaml")
	if err := viperObj.ReadInConfig(); err != nil {
		os.Create(filename)
	}
	viperObj.Set(key, value)
	return viperObj.WriteConfig()
}

func ConvertLinkArrayToMap(links []dal.Link) map[string]map[string]int {
	linkMap := make(map[string]map[string]int)
	for _, link := range links {
		if _, ok := linkMap[link.SourceNodeName]; !ok {
			linkMap[link.SourceNodeName] = make(map[string]int)
		}
		linkMap[link.SourceNodeName][link.TargetNodeName] = link.Bandwidth
	}
	return linkMap
}
