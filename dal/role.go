package dal

import "log"

//用于管理角色

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
