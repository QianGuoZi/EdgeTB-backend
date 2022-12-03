package handler

import (
	"EdgeTB-backend/dal"
	"EdgeTB-backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserInfo struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type UserData struct {
	Email string `form:"email" json:"email"`
}

type UserPwd struct {
	Password string `form:"password" json:"password"`
}

type Data struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expiredAt"`
}

// Login 用户登录
func Login(c *gin.Context) {
	// 获取用户名、密码
	var user UserInfo
	err := c.ShouldBind(&user)
	fmt.Println("Login传入的user信息", user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}
	// 校验用户名和密码是否正确
	// 生成Token
	token, expiredAt, err1 := service.Login(user.Username, user.Password)
	str := expiredAt.String()
	returnData := Data{token, str[0:10]}
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登陆成功",
		"data":    returnData,
	})
	return
}

// Register 用户注册
func Register(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	log.Printf("[Register] user=%+v", user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}

	_, err1 := service.Register(user.Username, user.Password)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "注册失败",
		})
	} else {
		token, times, err := service.GenerateToken(user.Username, user.Password)
		returnData := Data{token, times.String()}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "注册成功，自动登陆失败，请手动登陆",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "注册成功",
			"data":    returnData,
		})
	}
}

// GetUsername 测试用，解析token返回username
func GetUsername(c *gin.Context) {
	result, err := service.GetUsername(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": result,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功退出登录",
	})
}

// GetUserInfo 获取 token 返回用户信息
func GetUserInfo(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 根据用户名获取用户信息
	user, err := dal.GetUserInfoByName(username)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "获取信息成功",
	})
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取用户其他信息
	var data UserData
	err = c.ShouldBind(&data)
	log.Printf("[UpdateUserInfo] username=%s data=%+v", username, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户数据格式有误",
		})
		return
	}

	// 更新用户信息
	user := &dal.User{
		UserName: username,
		Email:    data.Email,
	}
	err = dal.UpdateUser(user)

	if err != nil {
		log.Printf("mysql update user failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "更新用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "用户信息修改成功",
	})
}

// UpdateUserPwd 更新用户密码
func UpdateUserPwd(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取用户新密码
	var userPwd UserPwd
	err = c.ShouldBind(&userPwd)
	pwd := userPwd.Password
	log.Printf("[UpdateUserPwd] password=%s", userPwd.Password)
	log.Printf("[UpdateUserPwd] username=%s", username)
	if pwd == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户密码不得为空",
		})
		return
	}

	userInfo, err1 := dal.GetUserInfoByName(username)
	if err1 != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 更新用户信息
	pwd, err = service.EncodePassword(username, pwd, userInfo.Salt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user := &dal.User{
		UserName: username,
		Pwd:      pwd,
	}
	err = dal.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "更新用户密码失败",
		})
		return
	}
	// 隐藏用户密码
	user.Pwd = ""
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		//"data":    user,
		"message": "修改密码成功",
	})
}
