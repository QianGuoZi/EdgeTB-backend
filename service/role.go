package service

import (
	"EdgeTB-backend/dal"
	"errors"
	"log"
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

type UploadedFile struct {
	FileName string `json:"fileName"`
	Size     int64  `json:"size"` // 单位字节
	URL      string `json:"url"`
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

// AddRole 添加角色
func AddRole(addRoleRequest RoleStruct, username string) error {
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddRole] 服务获取用户id失败")
		return errors.New("服务获取用户id失败")
	}
	check := dal.CheckRoleExist(addRoleRequest.Name, addRoleRequest.Image.Name, addRoleRequest.ImageSource)
	if check != nil {
		log.Printf("[AddRole] 服务添加角色失败，角色重复或镜像重复")
		return errors.New("角色重复或镜像重复")
	}
	//code
	codeId, err1 := AddRoleCode(addRoleRequest)
	if err1 != nil {
		return errors.New("添加角色代码信息失败")
	}
	//pyDep
	pyDepId, err2 := AddRolePyDep(addRoleRequest)
	if err2 != nil {
		return errors.New("添加角色py依赖信息失败")
	}
	//images
	imageId, err3 := AddRoleImage(addRoleRequest)
	if err3 != nil {
		return errors.New("添加角色image依赖信息失败")
	}
	//role
	roleId, err4 := AddRoleInfo(addRoleRequest, codeId, pyDepId, imageId, userId)
	if err4 != nil {
		return errors.New("添加角色信息失败")
	}
	//outputItem
	err5 := AddRoleOutputItem(addRoleRequest.OutputItems, roleId)
	if err5 != nil {
		return errors.New("添加角色输出信息失败")
	}
	return nil
}

// AddRoleCode 添加角色Code部分
func AddRoleCode(addRoleRequest RoleStruct) (int64, error) {
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
		codeId, err := dal.AddRoleCode(code)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色代码信息失败")
			return 0, errors.New("添加角色代码信息失败")
		}
		log.Printf("[AddRole] 服务添加角色codeId:%+v", codeId)
		return codeId, nil
	}
	return 0, nil
}

// AddRolePyDep 添加角色pyDep部分
func AddRolePyDep(addRoleRequest RoleStruct) (int64, error) {
	//pyDep:upload功能
	if addRoleRequest.PyDepSource == "upload" {
		rolePyDep := addRoleRequest.PyDep
		var pyDep dal.PyDev
		pyDep.PyDevSource = addRoleRequest.PyDepSource
		pyDep.PyDevPackages = *rolePyDep.Packages
		pyDepId, err := dal.AddRolePyDep(pyDep)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色py依赖信息失败")
			return 0, errors.New("添加角色py依赖信息失败")
		}
		log.Printf("[AddRole] 服务添加角色pyDepId:%+v", pyDepId)
		return pyDepId, nil
	}
	return 0, nil
}

// AddRoleImage 添加角色image部分
func AddRoleImage(addRoleRequest RoleStruct) (int64, error) {
	//image:platform功能
	if addRoleRequest.ImageSource == "platform" {
		roleImage := addRoleRequest.Image
		var image dal.Image
		image.ImageSource = addRoleRequest.ImageSource
		image.ImageName = roleImage.Name
		imageId, err := dal.AddRoleImage(image)
		if err != nil {
			log.Printf("[AddRole] 服务添加角色image信息失败")
			return 0, errors.New("添加角色image信息失败")
		}
		log.Printf("[AddRole] 服务添加角色imageId:%+v", imageId)
		return imageId, nil
	}
	return 0, nil
}

// AddRoleInfo 添加角色信息部分
func AddRoleInfo(addRoleRequest RoleStruct, codeId int64, pyDepId int64, imageId int64, userId int64) (int64, error) {
	var role dal.Role
	role.RoleName = addRoleRequest.Name
	role.Description = *addRoleRequest.Description
	role.PyVersion = addRoleRequest.PyVersion
	role.CodeId = codeId
	role.WorkDir = *addRoleRequest.WorkDir
	role.RunCommand = addRoleRequest.RunCommand
	role.PyDevId = pyDepId
	role.ImageId = imageId
	role.ImageName = addRoleRequest.Image.Name
	role.UserId = userId
	roleId, err := dal.AddRole(role)
	if err != nil {
		log.Printf("[AddRole] 服务添加角色信息失败")
		return 0, errors.New("添加角色信息失败")
	}
	log.Printf("[AddRole] 服务添加角色roleId:%+v", roleId)
	return roleId, nil
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
func GetAllRole(username string) ([]RoleListResponse, error) {
	var list []RoleListResponse
	//通过username获取id
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetAllRole] 服务获取用户id失败")
		return list, errors.New("服务获取用户id失败")
	}
	//获取name, description, pyVersion, imageId
	roleList, err1 := dal.GetAllRole(userId)
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
	//通过username获取id
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
	roleResult.PyDepSource = pyDepInfo.PyDevSource
	if pyDepInfo.PyDevSource == "upload" || pyDepInfo.PyDevSource == "manual" {
		roleResult.PyDep.Packages = new(string)
		*roleResult.PyDep.Packages = pyDepInfo.PyDevPackages
	} else if pyDepInfo.PyDevSource == "git" {
		roleResult.PyDep.Git = new(GitRepository)
		roleResult.PyDep.Git.Filepath = pyDepInfo.PyDevGitFilepath
		roleResult.PyDep.Git.URL = pyDepInfo.PyDevGitUrl
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
