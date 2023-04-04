package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"fmt"
	"log"
	"os"
)

type RoleStruct struct {
	Name        string       `json:"name"`                  // 角色名称，唯一
	Description *string      `json:"description,omitempty"` // 描述
	PyVersion   string       `json:"pyVersion"`             // python版本
	CodeSource  string       `json:"codeSource"`            // 代码来源
	Code        Code         `json:"code"`
	WorkDir     *string      `json:"workDir,omitempty"` // 工作目录，默认为代码根目录
	RunCommand  string       `json:"runCommand"`        // 启动指令
	PyDepSource string       `json:"pyDepSource"`       // py依赖来源
	PyDep       PyDep        `json:"pyDep"`
	ImageSource string       `json:"imageSource"` // 镜像来源
	Image       Image        `json:"image"`
	OutputItems []OutputItem `json:"outputItems"` // 输出项
}

type RoleUpdateRequest struct {
	Description *string      `json:"description,omitempty"` // 描述
	PyVersion   string       `json:"pyVersion"`             // python版本
	CodeSource  string       `json:"codeSource"`            // 代码来源
	Code        Code         `json:"code"`
	WorkDir     *string      `json:"workDir,omitempty"` // 工作目录，默认为代码根目录
	RunCommand  string       `json:"runCommand"`        // 启动指令
	PyDepSource string       `json:"pyDepSource"`       // py依赖来源
	PyDep       PyDep        `json:"pyDep"`
	ImageSource string       `json:"imageSource"` // 镜像来源
	Image       Image        `json:"image"`
	OutputItems []OutputItem `json:"outputItems"` // 输出项
}

type GitRepository struct {
	Filepath string `json:"filepath"` // 文件夹相对路径，dockerfile所在的文件夹路径
	URL      string `json:"url"`      // git仓库url
}

type Code struct {
	File   *UploadedFile `json:"file,omitempty"`   // 上传的代码文件，仅当代码来源为upload时具有
	GitURL *string       `json:"gitUrl,omitempty"` // 仅当代码来源为git时具有
}

type PyDep struct {
	Git      *GitRepository `json:"git,omitempty"`      // 包含req.txt文件的git仓库，仅当为git时具有
	Packages *string        `json:"packages,omitempty"` // 依赖库列表字符串，仅当为upload和manual时具有
}

type Image struct {
	Name       string         `json:"name"`                 // 镜像名称
	Dockerfile *UploadedFile  `json:"dockerfile,omitempty"` // dockerfile文件，仅当为dockerfile时具有
	Archive    *UploadedFile  `json:"archive,omitempty"`    // 包含dockerfile的压缩包文件，仅当为archive时具有
	Git        *GitRepository `json:"git,omitempty"`        // git仓库，仅当为git时具有
}

type OutputItem struct {
	Name string `json:"name"` // 输出项名称
	Path string `json:"path"` // 输出路径
}

type RoleListResponse struct {
	Name        string `json:"name"`                  // 角色名称，唯一
	Description string `json:"description,omitempty"` // 描述
	PyVersion   string `json:"pyVersion"`             // python版本
	ImageName   string `json:"imageName"`             // 镜像名称
}

// UploadRoleCodeFile 传入文件名、路径，获取类型、大小
func UploadRoleCodeFile(filePath string) (string, string, int, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[UploadRoleFile] 服务获取上传文件信息失败")
		return "", "", 0, errors.New("服务获取上传文件信息失败")
	}
	fmt.Println("name:", fi.Name())
	fmt.Println("size:", fi.Size())
	return filePath, fi.Name(), int(fi.Size()), nil
}

// AddRole 添加角色
func AddRole(addRoleRequest RoleStruct, username, projectName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddRole] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	check := dal.CheckRoleExist(addRoleRequest.Name)
	if check != nil {
		log.Printf("[AddRole] 服务添加角色失败，角色重复")
		return errors.New("角色重复")
	}
	projectId, err := dal.GetProjectId(projectName, userId)
	if err != nil {
		log.Printf("[AddRole] 服务获取项目id失败")
		return errors.New("服务获取项目id失败")
	}
	//role
	roleId, err1 := AddRoleInfo(addRoleRequest, userId, projectId)
	if err1 != nil {
		return errors.New("添加角色信息失败")
	}
	//code
	err = AddRoleCode(addRoleRequest, roleId)
	if err != nil {
		return errors.New("添加角色代码信息失败")
	}
	//pyDep
	err = AddRolePyDep(addRoleRequest, roleId)
	if err != nil {
		return errors.New("添加角色py依赖信息失败")
	}
	//images
	err = AddRoleImage(addRoleRequest, roleId)
	if err != nil {
		return errors.New("添加角色image依赖信息失败")
	}
	//outputItem
	err = AddRoleOutputItem(addRoleRequest.OutputItems, roleId)
	if err != nil {
		return errors.New("添加角色输出信息失败")
	}
	return nil
}

// AddRoleInfo 添加角色信息部分
func AddRoleInfo(addRoleRequest RoleStruct, userId, projectId int64) (int64, error) {
	var role dal.Role
	role.RoleName = addRoleRequest.Name
	role.Description = *addRoleRequest.Description
	role.PyVersion = addRoleRequest.PyVersion
	role.WorkDir = *addRoleRequest.WorkDir
	role.RunCommand = addRoleRequest.RunCommand
	role.ImageName = addRoleRequest.Image.Name
	role.ProjectId = projectId
	role.UserId = userId
	roleId, err := dal.AddRole(role)
	if err != nil {
		log.Printf("[AddRole] 服务添加角色信息失败")
		return 0, errors.New("添加角色信息失败")
	}
	log.Printf("[AddRole] 服务添加角色roleId:%+v", roleId)
	return roleId, nil
}

// AddRoleCode 添加角色Code部分
func AddRoleCode(addRoleRequest RoleStruct, roleId int64) error {
	//code:upload功能
	if addRoleRequest.CodeSource == "upload" {
		//拿到文件信息，传入数据库
		roleCode := addRoleRequest.Code
		file := roleCode.File
		var code dal.Code
		code.CodeSource = addRoleRequest.CodeSource
		code.CodeFileName = file.FileName
		code.CodeFileSize = file.Size
		code.CodeFileUrl = file.URL
		code.RoleId = roleId
		codeId, err := dal.AddRoleCode(code)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色代码信息失败")
			return errors.New("添加角色代码信息失败")
		}
		log.Printf("[AddRole] 服务添加角色codeId:%+v", codeId)
		return nil
	}
	return nil
}

// AddRolePyDep 添加角色pyDep部分
func AddRolePyDep(addRoleRequest RoleStruct, roleId int64) error {
	//pyDep:upload功能
	if addRoleRequest.PyDepSource == "upload" || addRoleRequest.PyDepSource == "manual" {
		rolePyDep := addRoleRequest.PyDep
		var pyDep dal.PyDep
		pyDep.PyDepSource = addRoleRequest.PyDepSource
		pyDep.PyDepPackages = *rolePyDep.Packages
		pyDep.RoleId = roleId
		pyDepId, err := dal.AddRolePyDep(pyDep)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色py依赖信息失败")
			return errors.New("添加角色py依赖信息失败")
		}
		log.Printf("[AddRole] 服务添加角色pyDepId:%+v", pyDepId)
		return nil
	}
	return nil
}

// AddRoleImage 添加角色image部分
func AddRoleImage(addRoleRequest RoleStruct, roleId int64) error {
	//image:platform功能
	if addRoleRequest.ImageSource == "platform" {
		roleImage := addRoleRequest.Image
		var image dal.Image
		image.ImageSource = addRoleRequest.ImageSource
		image.ImageName = roleImage.Name
		image.RoleId = roleId
		imageId, err := dal.AddRoleImage(image)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色image信息失败")
			return errors.New("添加角色image信息失败")
		}
		log.Printf("[AddRole] 服务添加角色imageId:%+v", imageId)
		return nil
	}
	return nil
}

// AddRoleOutputItem 添加角色outputItem部分
func AddRoleOutputItem(outputItems []OutputItem, roleId int64) error {
	listLen := len(outputItems)
	itemList := make([]dal.OutputItem, listLen)
	for i := 0; i < listLen; i++ {
		itemList[i].OutputName = outputItems[i].Name
		itemList[i].OutputPath = outputItems[i].Path
		itemList[i].RoleId = roleId
	}
	err := dal.AddRoleOutputItem(itemList)
	if err != nil {
		return errors.New("添加角色outputItem失败")
	}
	return nil
}

// GetAllRole 获取用户所有角色基本信息
func GetAllRole(username string, projectName string) ([]RoleListResponse, error) {
	var list []RoleListResponse
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetAllRole] 服务获取用户id失败")
		return list, errors.New("服务获取用户id失败")
	}
	//获取name, description, pyVersion, imageId
	var roleList []dal.Role
	var err1 error
	if projectName == "" {
		roleList, err1 = dal.GetAllRole(userId)
	} else {
		var projectId int64
		projectId, err1 = dal.GetProjectId(projectName, userId)
		if err1 != nil {
			log.Printf("[GetAllRole] 服务获取项目id失败")
			return list, errors.New("服务获取项目id失败")
		}
		roleList, err1 = dal.GetAllRoleByProjectId(projectId)
	}
	if err1 != nil {
		log.Printf("[GetAllRole] 服务获取用户角色列表失败")
		return list, errors.New("服务获取用户角色列表失败")
	}
	listLen := len(roleList)
	responseList := make([]RoleListResponse, listLen)
	for i := 0; i < listLen; i++ {
		responseList[i].Name = roleList[i].RoleName
		responseList[i].Description = roleList[i].Description
		responseList[i].PyVersion = roleList[i].PyVersion
		responseList[i].ImageName = roleList[i].ImageName
	}
	return responseList, nil
}

// GetRoleDetail 获取角色详细信息
func GetRoleDetail(username, roleName string) (RoleStruct, error) {
	var roleResult *RoleStruct
	roleResult = new(RoleStruct)
	//通过username获取userId
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetAllRole] 服务获取用户id失败")
		return *roleResult, errors.New("服务获取用户id失败")
	}
	//查user和role，返回role
	roleInfo, err1 := dal.GetRoleInfo(userId, roleName)
	if err1 != nil {
		log.Printf("[GetAllRole] 服务获取角色信息失败")
		return *roleResult, errors.New("服务获取角色信息失败")
	}
	roleId := roleInfo.Id
	roleResult.Name = roleInfo.RoleName
	if roleInfo.Description != "" {
		roleResult.Description = new(string)
		*roleResult.Description = roleInfo.Description
	}
	if roleInfo.WorkDir != "" {
		roleResult.WorkDir = new(string)
		*roleResult.WorkDir = roleInfo.WorkDir
	}
	roleResult.PyVersion = roleInfo.PyVersion
	roleResult.RunCommand = roleInfo.RunCommand

	//code
	codeInfo, err2 := dal.GetRoleCode(roleId)
	if err2 != nil {
		log.Printf("[GetAllRole] 服务获取角色code信息失败")
		return *roleResult, errors.New("服务获取角色code信息失败")
	}
	roleResult.CodeSource = codeInfo.CodeSource
	if codeInfo.CodeSource == "upload" {
		roleResult.Code.File = new(UploadedFile)
		roleResult.Code.File.FileName = codeInfo.CodeFileName
		roleResult.Code.File.Size = codeInfo.CodeFileSize
		roleResult.Code.File.URL = codeInfo.CodeFileUrl
	} else if codeInfo.CodeSource == "git" {
		roleResult.Code.GitURL = new(string)
		*roleResult.Code.GitURL = codeInfo.CodeGitUrl
	}

	//pyDep
	pyDepInfo, err3 := dal.GetRolePyDep(roleId)
	if err3 != nil {
		log.Printf("[GetAllRole] 服务获取角色pyDep信息失败")
		return *roleResult, errors.New("服务获取角色pyDep信息失败")
	}
	roleResult.PyDepSource = pyDepInfo.PyDepSource
	if pyDepInfo.PyDepSource == "upload" || pyDepInfo.PyDepSource == "manual" {
		roleResult.PyDep.Packages = new(string)
		*roleResult.PyDep.Packages = pyDepInfo.PyDepPackages
	} else if pyDepInfo.PyDepSource == "git" {
		roleResult.PyDep.Git = new(GitRepository)
		roleResult.PyDep.Git.Filepath = pyDepInfo.PyDepGitFilepath
		roleResult.PyDep.Git.URL = pyDepInfo.PyDepGitUrl
	}

	//image
	imageInfo, err4 := dal.GetRoleImage(roleId)
	if err4 != nil {
		log.Printf("[GetAllRole] 服务获取角色image信息失败")
		return *roleResult, errors.New("服务获取角色image信息失败")
	}
	roleResult.ImageSource = imageInfo.ImageSource
	roleResult.Image.Name = imageInfo.ImageName
	if imageInfo.ImageSource == "git" {
		roleResult.Image.Git = new(GitRepository)
		roleResult.Image.Git.Filepath = imageInfo.ImageGitFilepath
		roleResult.Image.Git.URL = imageInfo.ImageGitUrl
	} else if imageInfo.ImageSource == "uploadArchive" {
		roleResult.Image.Archive = new(UploadedFile)
		roleResult.Image.Archive.FileName = imageInfo.ImageArchiveName
		roleResult.Image.Archive.Size = imageInfo.ImageArchiveSize
		roleResult.Image.Archive.URL = imageInfo.ImageArchiveUrl
	} else if imageInfo.ImageSource == "uploadDockerfile" {
		roleResult.Image.Dockerfile = new(UploadedFile)
		roleResult.Image.Dockerfile.FileName = imageInfo.ImageDockerfileName
		roleResult.Image.Dockerfile.Size = imageInfo.ImageDockerfileSize
		roleResult.Image.Dockerfile.URL = imageInfo.ImageDockerfileUrl
	}

	//outputItem
	outputItemInfo, listLen, err5 := dal.GetRoleOutputItem(roleId)
	if err5 != nil {
		log.Printf("[GetAllRole] 服务获取角色outputItem信息失败")
		return *roleResult, errors.New("服务获取角色outputItem信息失败")
	}
	outputItemList := make([]OutputItem, listLen)
	for i := 0; i < listLen; i++ {
		outputItemList[i].Name = outputItemInfo[i].OutputName
		outputItemList[i].Path = outputItemInfo[i].OutputPath
	}
	roleResult.OutputItems = outputItemList

	return *roleResult, nil
}

// UpdateRole 更新角色信息
func UpdateRole(roleName string, roleUpdateInfo RoleUpdateRequest, username string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[UpdateRole] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	//通过roleName获取roleId
	roleId, err1 := dal.GetRoleId(userId, roleName)
	if err1 != nil {
		log.Printf("[UpdateRole] 服务修改角色信息失败，用户无权修改角色")
		return errors.New("用户无权修改角色")
	}

	//role：判断userId和roleId后，直接修改description，py_version，work_dir，run_command
	err = dal.UpdateRoleInfo(roleId, *roleUpdateInfo.Description, roleUpdateInfo.PyVersion,
		*roleUpdateInfo.WorkDir, roleUpdateInfo.RunCommand)
	if err != nil {
		log.Printf("[UpdateRole] 服务修改角色信息失败")
		return errors.New("服务修改角色信息失败")
	}
	//code：直接修改
	err = dal.UpdateRoleCode(roleId, roleUpdateInfo.Code.File.Size, roleUpdateInfo.CodeSource, roleUpdateInfo.Code.File.FileName,
		roleUpdateInfo.Code.File.URL, *roleUpdateInfo.Code.GitURL)
	if err != nil {
		log.Printf("[UpdateRole] 服务修改角色code信息失败")
		return errors.New("服务修改角色code信息失败")
	}
	//pyDep：直接修改
	err = dal.UpdateRolePyDep(roleId, roleUpdateInfo.PyDepSource, *roleUpdateInfo.PyDep.Packages,
		roleUpdateInfo.PyDep.Git.URL, roleUpdateInfo.PyDep.Git.Filepath)
	if err != nil {
		log.Printf("[UpdateRole] 服务修改角色pyDep信息失败")
		return errors.New("服务修改角色pyDep信息失败")
	}
	//images：直接修改+修改role中的image_name
	err = dal.UpdateRoleImage(roleId, roleUpdateInfo.ImageSource, roleUpdateInfo.Image.Name,
		roleUpdateInfo.Image.Dockerfile.URL, roleUpdateInfo.Image.Dockerfile.FileName,
		roleUpdateInfo.Image.Archive.URL, roleUpdateInfo.Image.Archive.FileName,
		roleUpdateInfo.Image.Git.URL, roleUpdateInfo.Image.Git.Filepath,
		roleUpdateInfo.Image.Dockerfile.Size, roleUpdateInfo.Image.Archive.Size)
	if err != nil {
		log.Printf("[UpdateRole] 服务修改角色image信息失败")
		return errors.New("服务修改角色image信息失败")
	}
	//outputItem：删除后添加
	err = dal.DeleteRoleOutputItem(roleId)
	if err != nil {
		log.Printf("[UpdateRole] 服务删除角色outputItem信息失败")
		return errors.New("服务修改角色outputItem信息失败")
	}
	err = AddRoleOutputItem(roleUpdateInfo.OutputItems, roleId)
	if err != nil {
		log.Printf("[UpdateRole] 服务添加角色outputItem信息失败")
		return errors.New("服务修改角色outputItem信息失败")
	}
	return nil
}

// DeleteRole 删除角色信息
func DeleteRole(username, roleName string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[DeleteRole] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	//通过roleName获取roleId
	roleId, err1 := dal.GetRoleId(userId, roleName)
	if err1 != nil {
		log.Printf("[DeleteRole] 服务删除角色信息失败，用户无权修改角色")
		return errors.New("用户无权删除角色")
	}
	//code
	err = dal.DeleteRoleCode(roleId)
	if err != nil {
		log.Printf("[DeleteRole] 服务删除角色code信息失败")
		return errors.New("服务删除角色code信息失败")
	}
	//pyDep
	err = dal.DeleteRolePyDep(roleId)
	if err != nil {
		log.Printf("[DeleteRole] 服务删除角色pyDep信息失败")
		return errors.New("服务删除角色pyDep信息失败")
	}
	//image
	err = dal.DeleteRoleImage(roleId)
	if err != nil {
		log.Printf("[DeleteRole] 服务删除角色image信息失败")
		return errors.New("服务删除角色image信息失败")
	}
	//outputItem
	err = dal.DeleteRoleOutputItem(roleId)
	if err != nil {
		log.Printf("[DeleteRole] 服务删除角色outputItem信息失败")
		return errors.New("服务修改角色outputItem信息失败")
	}
	//role
	err = dal.DeleteRoleInfo(roleId)
	if err != nil {
		log.Printf("[DeleteRole] 服务删除角色信息失败")
		return errors.New("服务修改角色信息失败")
	}
	return nil
}
