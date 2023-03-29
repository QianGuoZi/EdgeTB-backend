package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

type RoleDetailRequest struct {
	Name string `json:"name"` // 角色名称，唯一
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
	//获取项目名称
	projectName, _ := c.GetQuery("project")
	log.Printf("[AddRole] projectName=%+v", projectName)
	//获取数据
	var newRole service.RoleStruct
	err1 := c.ShouldBind(&newRole)
	log.Printf("[AddRole] newRole=%+v", newRole)
	newRole.Description = new(string)
	newRole.WorkDir = new(string)
	newRole.Code.File = new(service.UploadedFile)
	newRole.Code.GitURL = new(string)
	newRole.PyDep.Packages = new(string)
	newRole.PyDep.Git = new(service.GitRepository)
	newRole.Image.Git = new(service.GitRepository)
	newRole.Image.Archive = new(service.UploadedFile)
	newRole.Image.Dockerfile = new(service.UploadedFile)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "新增角色数据格式有误",
		})
		return
	}
	//传给Service层处理
	err = service.AddRole(newRole, username, projectName)
	//返回成功或失败消息
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "新增角色失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "新增角色成功",
	})
	return
}

// AllRole 列出角色
func AllRole(c *gin.Context) {
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
	//service层获取用户的角色列表数据
	returnData, err1 := service.GetAllRole(username)
	//返回角色数据
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取角色列表失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取角色列表成功",
		"data":    returnData,
	})
	return
}

// RoleDetail 查看角色详情
func RoleDetail(c *gin.Context) {
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
	//获取角色name
	roleName := c.Param("name")
	log.Printf("[RoleDetail] roleName=%+v", roleName)
	if roleName == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "角色名数据格式有误",
		})
		return
	}
	//service层获取角色详情
	returnData, err2 := service.GetRoleDetail(username, roleName)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取角色详情失败",
		})
		return
	}
	//返回角色详情
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取角色详情成功",
		"data":    returnData,
	})
	return
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
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
	//获取角色name
	roleName := c.Param("name")
	log.Printf("[UpdateRole] roleName=%+v", roleName)
	if roleName == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "角色名数据格式有误",
		})
		return
	}
	//获取修改的数据
	var roleUpdateInfo service.RoleUpdateRequest
	roleUpdateInfo.Description = new(string)
	roleUpdateInfo.WorkDir = new(string)
	roleUpdateInfo.Code.File = new(service.UploadedFile)
	roleUpdateInfo.Code.GitURL = new(string)
	roleUpdateInfo.PyDep.Packages = new(string)
	roleUpdateInfo.PyDep.Git = new(service.GitRepository)
	roleUpdateInfo.Image.Git = new(service.GitRepository)
	roleUpdateInfo.Image.Archive = new(service.UploadedFile)
	roleUpdateInfo.Image.Dockerfile = new(service.UploadedFile)
	err1 := c.ShouldBind(&roleUpdateInfo)
	log.Printf("[UpdateRole] roleUpdate=%+v", roleUpdateInfo)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "更新角色数据格式有误",
		})
		return
	}
	//service层处理修改
	err = service.UpdateRole(roleName, roleUpdateInfo, username)
	//返回成功或失败消息
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "角色修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色修改成功",
	})
	return

}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	//获取username
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
	//获取角色name
	roleName := c.Param("name")
	log.Printf("[RoleDetail] roleName=%+v", roleName)
	if roleName == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "角色名数据格式有误",
		})
		return
	}
	//service层处理删除
	err = service.DeleteRole(username, roleName)
	//返回成功或失败
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "角色删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色删除成功",
	})
	return
}

// UploadRoleCode 上传本地代码文件
func UploadRoleCode(c *gin.Context) {
	//上传文件
	//上传到本地后添加新文件信息
	//获取用户名称
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token有误",
		})
		return
	}

	//处理文件上传
	//目录
	var saveDir = ""
	//名称
	var saveName = ""
	//完整路径
	var savePath = ""
	//获取文件
	file, errFile := c.FormFile("file")
	//处理获取文件错误
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择文件",
		})
		return
	}
	//获取上传文件的文件名和后缀(类型)
	uploadFileNameWithSuffix := path.Base(file.Filename) //得到完整文件名
	uploadFileType := path.Ext(uploadFileNameWithSuffix) //得到.txt等
	log.Printf("[UploadRoleCode] fileName=%+v", file.Filename)
	log.Printf("[UploadRoleCode] uploadFileType=%+v", uploadFileType)
	//判断后缀是否合法
	checkResult := CheckZipFile(uploadFileType)
	log.Printf("[UploadRoleCode] checkResult=%+v", checkResult)
	if checkResult != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件类型有误",
		})
		return
	}
	//文件保存目录
	saveDir = "./RoleFile/CodeFile"
	//保存的文件名称
	dayStr := time.Now().String()
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	saveName = username + "_" + dayStr[0:10] + "_" + timeStamp + "_" + file.Filename
	//文件保存的路径
	savePath = saveDir + "/" + saveName
	log.Printf("[UploadRoleCode] savePath=%+v", savePath)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件上传失败",
		})
		return
	}
	//没有错误的情况下
	//返回文件path,size和fileName
	filePath, fileName, size, err1 := service.UploadRoleCodeFile(savePath)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件上传成功，但文件信息获取失败",
		})
		return
	}
	returnData := FileReturn{filePath, fileName, size}
	log.Printf("[UploadRoleCode] returnData=%+v", returnData)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data":    returnData,
	})
	return
}
