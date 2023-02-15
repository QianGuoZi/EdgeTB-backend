package handler

import (
	"EdgeTB-backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
	var newRole service.AddRoleRequest
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
	err = service.AddRole(newRole, username)
	//返回成功或失败消息
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "新增角色失败",
		})
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
