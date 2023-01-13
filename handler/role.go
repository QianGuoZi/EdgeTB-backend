package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AddRoleRequest struct {
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

// AddRole 创建角色
func AddRole(c *gin.Context) {
	//获取用户username
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}
	log.Printf("[GetUserInfo] success username=%+v", username)
	//获取数据
	var newRole AddRoleRequest
	err1 := c.ShouldBind(&newRole)
	log.Printf("[AddRole] newRole=%+v", newRole)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "新增角色数据格式有误",
		})
		return
	}
	//传给Service层处理

	//返回成功或失败消息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取公共数据集列表成功",
		"data":    newRole,
	})
	return
}

// AllRole 列出角色
func AllRole(c *gin.Context) {
	//获取用户username
	//service层获取用户的角色列表数据
	//返回角色数据
}

// RoleDetail 查看角色详情
func RoleDetail(c *gin.Context) {
	//获取用户username
	//获取角色name
	//service层获取角色详情
	//返回角色详情
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	//获取username
	//获取角色name
	//获取修改的数据
	//service层处理修改
	//返回成功或失败
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	//获取username
	//获取角色name
	//service层处理删除
	//返回成功或失败
}

// UploadRoleCode 上传本地代码文件
func UploadRoleCode(c *gin.Context) {
	//获取username
	//获取文件
	//service层处理文件信息
	//返回文件信息
}
