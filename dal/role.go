package dal

import (
	"errors"
	"log"
)

//用于管理角色

// CheckRoleExist 检查是否已存在image和role
func CheckRoleExist(roleName, imageName, imageSource string) error {
	var checkImage Image
	check := DB.Model(&Image{}).Where("image_name = ? && image_source = ?", imageName, imageSource).First(&checkImage)
	if check.Error == nil {
		log.Printf("[CheckRoleExist] image信息已存在")
		return errors.New("image信息已存在")
	}
	var checkRole Role
	check = DB.Model(&Role{}).Where("role_name = ? ", roleName).First(&checkRole)
	if check.Error == nil {
		log.Printf("[CheckRoleExist] role信息已存在")
		return errors.New("role信息已存在")
	}
	return nil
}

// AddRoleCode 添加角色code信息
func AddRoleCode(code Code) (int64, error) {
	result := DB.Model(&Code{}).Create(&code)
	if result.Error != nil {
		log.Printf("[AddRoleCode] 数据库创建角色code信息失败")
		return 0, result.Error
	}
	return code.Id, nil
}

// AddRolePyDep 添加角色pyDep信息
func AddRolePyDep(pyDev PyDev) (int64, error) {
	result := DB.Model(&PyDev{}).Create(&pyDev)
	if result.Error != nil {
		log.Printf("[AddRolePyDep] 数据库创建角色pyDep信息失败")
		return 0, result.Error
	}
	return pyDev.Id, nil
}

// AddRoleImage 添加角色image信息
func AddRoleImage(image Image) (int64, error) {
	result := DB.Model(&Image{}).Create(&image)
	if result.Error != nil {
		log.Printf("[AddRoleImage] 数据库创建角色image信息失败")
		return 0, result.Error
	}
	return image.Id, nil
}

// AddRole 添加角色信息
func AddRole(role Role) (int64, error) {
	result := DB.Model(&Role{}).Create(&role)
	if result.Error != nil {
		log.Printf("[AddRole] 数据库创建角色信息失败")
		return 0, result.Error
	}
	return role.Id, nil
}

// AddRoleOutputItem 添加outputItem信息
func AddRoleOutputItem(outputItemList []OutputItem) error {
	result := DB.Model(&OutputItem{}).Create(&outputItemList)
	if result.Error != nil {
		log.Printf("[AddRoleOutputItem] 数据库创建outputItem信息失败")
		return result.Error
	}
	return nil
}

// GetAllRole 数据库获取角色列表信息
func GetAllRole(userId int64) ([]Role, error) {
	var roles []Role
	result := DB.Model(&Role{}).Select("id", "role_name", "description", "py_version", "image_name").
		Where("user_id = ?", userId).Find(&roles)
	if result.Error != nil {
		log.Printf("[GetAllRoleInfo] 数据库获取角色列表信息失败")
		return roles, result.Error
	}
	listLen := len(roles)
	log.Printf("[GetAllRoleInfo] 数据库获取角色列表信息成功，长度为：%+v", listLen)
	return roles, nil
}

// GetRoleInfo 获取角色详细信息
func GetRoleInfo(userId int64, roleName string) (Role, error) {
	var roleInfo Role
	result := DB.Model(&Role{}).Where("user_id = ? && role_name = ?", userId, roleName).First(&roleInfo)
	if result.Error != nil {
		log.Printf("[GetRoleInfo] 数据库获取角色信息失败")
		return roleInfo, result.Error
	}
	log.Printf("[GetRoleInfo] 数据库获取角色信息成功")
	return roleInfo, nil
}

// GetRoleCode 获取角色code信息
func GetRoleCode(codeId int64) (Code, error) {
	var codeInfo Code
	result := DB.Model(&Code{}).Where("id = ?", codeId).First(&codeInfo)
	if result.Error != nil {
		log.Printf("[GetRoleCode] 数据库获取角色code信息失败")
		return codeInfo, result.Error
	}
	log.Printf("[GetRoleCode] 数据库获取角色code信息成功")
	return codeInfo, nil
}

// GetRolePyDep 获取角色pyDep信息
func GetRolePyDep(codeId int64) (PyDev, error) {
	var pyDepInfo PyDev
	result := DB.Model(&PyDev{}).Where("id = ?", codeId).First(&pyDepInfo)
	if result.Error != nil {
		log.Printf("[GetRolePyDep] 数据库获取角色pyDep信息失败")
		return pyDepInfo, result.Error
	}
	log.Printf("[GetRolePyDep] 数据库获取角色pyDep信息成功")
	return pyDepInfo, nil
}

// GetRoleImage 获取角色image信息
func GetRoleImage(codeId int64) (Image, error) {
	var imageInfo Image
	result := DB.Model(&Image{}).Where("id = ?", codeId).First(&imageInfo)
	if result.Error != nil {
		log.Printf("[GetRolePyDep] 数据库获取角色image信息失败")
		return imageInfo, result.Error
	}
	log.Printf("[GetRolePyDep] 数据库获取角色image信息成功")
	return imageInfo, nil
}

// GetRoleOutputItem 获取角色outputItem信息
func GetRoleOutputItem(roleId int64) ([]OutputItem, int, error) {
	var outputItemList []OutputItem
	result := DB.Model(&OutputItem{}).Where("role_id = ?", roleId).Find(&outputItemList)
	if result.Error != nil {
		log.Printf("[GetRoleOutputItem] 数据库获取角色outputItem信息失败")
		return outputItemList, 0, result.Error
	}
	listLen := len(outputItemList)
	log.Printf("[GetRoleOutputItem] 数据库获取角色outputItem信息成功，长度为%+v", listLen)
	return outputItemList, listLen, nil
}
