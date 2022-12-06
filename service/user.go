package service

import (
	"EdgeTB-backend/dal"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strings"
	"time"
)

//随机生成长度为4的盐
func randSalt() string {
	rand.Seed(time.Now().UnixNano())
	buf := strings.Builder{}
	for i := 0; i < 4; i++ {
		// 如果写byte会无法兼容mysql编码
		buf.WriteRune(rune(rand.Intn(256)))
	}
	return buf.String()
}

// Register 用userName和password注册用户
func Register(userName, password string) (id int64, err error) {
	//查看用户是否存在
	user, err := dal.CheckUser(userName)
	if err != nil {
		return 0, errors.New("该用户已存在")
	}

	//设置salt，并生成pwd
	user.Salt = randSalt()
	//pw := md5.New()
	//pw.Write([]byte(password))
	//password = hex.EncodeToString(pw.Sum(nil))
	//fmt.Println("md5:", password)
	pwd, err := EncodePassword(userName, password, user.Salt)
	if err != nil {
		return 0, err
	}

	user.Id = id
	user.UserName = userName
	user.Pwd = pwd

	log.Printf("[Register] database new user=%+v", user)
	returnId, err := dal.AddUser(user)
	if err != nil {
		return 0, err
	}

	return returnId, nil
}

//Login 使用userName password登陆，返回token和time
func Login(userName, password string) (string, time.Time, error) {
	result, err := dal.SearchUser(userName, password)
	if err != nil {
		log.Printf("[Login] dal.SearchUser执行失败")
		return "", time.Now(), errors.New("鉴权失败")
	} else if result == false {
		return "", time.Now(), errors.New("用户名或密码错误")
	}
	token, times, err := GenerateToken(userName)
	if err != nil {
		log.Printf("[Login] GenerateToken失败")
		return "", time.Now(), errors.New("生成token失败")
	}
	return token, times, nil
}

// GetUsername 基于JWT的认证中间件获取用户名
func GetUsername(c *gin.Context) (string, error) {
	// Token放在Header的Authorization中，并使用Bearer开头
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "auth为空", errors.New("auth为空")
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "请求头中auth格式有误", errors.New("请求头中auth格式有误")
	}
	// parts[1]是获取到的tokenString，用定义好的解析JWT的函数来解析它
	mc, err := ParseToken(parts[1])
	if err != nil {
		return "无效的token", errors.New("无效的token")
	}
	//返回userName
	return mc.Username, nil
}

// EncodePassword 将密码进行转换
func EncodePassword(userName, password, salt string) (string, error) {
	//设置salt，并生成pwd
	buf := bytes.Buffer{}
	buf.WriteString(userName)
	buf.WriteString(password)
	buf.WriteString(salt)
	pwd, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.MinCost)
	if err != nil {
		return "", errors.New("密码加盐失败")
	}
	return string(pwd), nil
}
