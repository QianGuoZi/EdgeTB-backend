package dal

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// CheckUser 查看用户是否已存在
func CheckUser(userName string) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ?", userName).First(&user)
	if user.Id != 0 {
		log.Printf("[CheckUser] 用户存在")
		return user, errors.New("用户存在")
	}
	return user, nil
}

// AddUser 向users表中添加新记录
func AddUser(user User) (int64, error) {
	result := DB.Model(&User{}).Create(&user)
	if result.Error != nil {
		log.Printf("[Register] 数据库创建用户记录失败")
		return 0, result.Error
	}
	return user.Id, nil
}

// SearchUser 判断userName password是否正确
func SearchUser(userName, password string) (bool, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ?", userName).First(&user)
	log.Printf("[SearchUser] search user=%+v", user)
	if user.UserName != userName {
		log.Printf("[SearchUser] 用户名错误")
		return false, errors.New("用户名错误")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(userName+password+user.Salt))
	if err != nil {
		log.Printf("[SearchUser] 密码错误")
		return false, errors.New("密码错误")
	}
	return true, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(id int64) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("id = ?", id).First(&user)
	if user.Id == 0 {
		return User{}, errors.New("无法找到该用户")
	}
	//屏蔽掉密码等信息
	user.Pwd = ""
	user.Salt = ""
	fmt.Println("getUserInfo", user)
	return user, nil
}

// GetUserInfoByName 获取用户信息
func GetUserInfoByName(name string) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ?", name).First(&user)
	if user.Id == 0 {
		return User{}, errors.New("无法找到该用户")
	}
	//屏蔽掉密码等信息
	user.Pwd = ""
	user.Salt = ""
	log.Printf("GetUserInfoByName usr=%+v", user)
	return user, nil
}

// UpdateUser 更改用户信息
func UpdateUser(user *User) error {
	info := make(map[string]interface{})
	if user.Email != "" {
		info["company"] = user.Email
	}
	if user.Pwd != "" {
		info["pwd"] = user.Pwd
	}
	log.Printf("[mysql UpdateUser] user=%+v data=%+v", user, info)
	result := DB.Model(&User{}).Where("user_name = ?", user.UserName).Updates(info)

	if result.Error != nil {
		log.Printf("mysql update user failed err=%+v", result.Error)
		return result.Error
	}
	return nil
}
