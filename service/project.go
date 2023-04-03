package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
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
	// Dataset    Dataset    `json:"dataset"`    // 数据集
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
	// //dataset 信息添加
	// //datasetSplitter
	// var datasetSplitterFile dal.File
	// datasetSplitterFile.Name = projectInfo.Dataset.Splitter.FileName
	// datasetSplitterFile.Url = projectInfo.Dataset.Splitter.URL
	// datasetSplitterFile.Size = projectInfo.Dataset.Splitter.Size
	// datasetSplitterId, err3 := dal.AddFile(datasetSplitterFile)
	// if err3 != nil {
	// 	log.Printf("[AddProjectDetail] 服务添加datasetSplitter文件信息失败")
	// 	return errors.New("服务添加datasetSplitter文件信息失败")
	// }
	// //存储4个id
	// err = dal.UpdateProjectInfo(projectId, managerId, structureId, int64(projectInfo.Dataset.ID), datasetSplitterId)
	err = dal.UpdateProjectInfo(projectId, managerId, structureId, 0, 0)
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

// StartProject 运行项目
func StartProject(username, projectName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetProjectDetail] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	projectId, err1 := dal.GetProjectId(projectName, userId)
	if err1 != nil {
		log.Printf("[GetProjectDetail] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	//创建yaml文件
	err = CreateConfigYaml(username, projectName)
	if err != nil {
		return err
	}
	//role文件
	//TODO:装模做样？
	//manager文件
	//文件位置./ProjectFile/Manager
	//数据集切分脚本文件
	//文件位置./ProjectFile/DatasetSplitter 执行python3 dataset_conf.py -d dataset.json
	cmd("python3 /home/qianguo/controller/dml_tool/dataset_conf.py -d /home/qianguo/controller/dml_tool/gl_dataset.json")
	//DML结构配置脚本文件
	//文件位置./ProjectFile/Structure 执行python3 gl_structure_conf.py -s structure.json
	cmd("python3 /home/qianguo/controller/dml_tool/gl_structure_conf.py -s /home/qianguo/controller/dml_tool/gl_structure.json")
	//运行controller
	go cmdRun("sudo PROJECT_ID=" + strconv.Itoa(int(projectId)) + "BACKEND_ADDR=127.0.0.1:3000 python3 /home/qianguo/controller/gl_run.py\n")
	//curl localhost:3333/conf/dataset
	time.Sleep(time.Second * 30)
	cmd("curl localhost:3333/conf/dataset")
	//curl localhost:3333/conf/structure
	cmd("curl localhost:3333/conf/structure")
	//log接收到tc finish后
	time.Sleep(time.Second * 60) //等tc
	//curl localhost:3333/start
	cmd("curl localhost:3333/start")
	return nil
}

// FinishProject 终止项目
func FinishProject(username, projectName string) error {
	//curl localhost:3333/finish
	cmdRun("curl localhost:3333/finish")
	log.Printf(username, projectName)
	return nil
}

func cmdRun(c string) {
	cmd := exec.Command("/bin/bash", "-c", c) //不加第一个第二个参数会报错

	//cmd.Stdout = os.Stdout // cmd.Stdout -> stdout  重定向到标准输出，逐行实时打印
	//cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	//也可以重定向文件 cmd.Stderr= fd (文件打开的描述符即可)

	stdout, _ := cmd.StdoutPipe() //创建输出管道
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	fmt.Println(cmd.Args) //查看当前执行命令
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
	return
}

func cmd(c string) {
	cmd := exec.Command("/bin/bash", "-c", c) //不加第一个第二个参数会报错

	//cmd.Stdout = os.Stdout // cmd.Stdout -> stdout  重定向到标准输出，逐行实时打印
	//cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	//也可以重定向文件 cmd.Stderr= fd (文件打开的描述符即可)

	stdout, _ := cmd.StdoutPipe() //创建输出管道
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	fmt.Println(cmd.Args) //查看当前执行命令

	cmdPid := cmd.Process.Pid //查看命令pid
	fmt.Println(cmdPid)

	result, _ := ioutil.ReadAll(stdout) // 读取输出结果
	resData := string(result)
	fmt.Println(resData)

	var res int
	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			fmt.Println("cmd exit status")
			res = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
		}
	}
	fmt.Println(res)
}
