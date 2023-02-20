package dal

import (
	"errors"
	"log"
)

//用于管理角色
//TODO:image的存在与唯一未判断

// GetRoleId 检查用户与角色并获取角色id
func GetRoleId(userId int64, roleName string) (int64, error) {
	role := Role{}
	DB.Model(&Role{}).Where("role_name = ? && user_id = ?", roleName, userId).First(&role)
	if role.Id == 0 {
		return 0, errors.New("无法找到该角色")
	}
	log.Printf("[GetRoleId] roleId=%+v", role.Id)
	return role.Id, nil
}

// CheckRoleExist 检查是否已存在role
func CheckRoleExist(roleName string) error {
	var checkRole Role
	check := DB.Model(&Role{}).Where("role_name = ? ", roleName).First(&checkRole)
	if check.Error == nil {
		log.Printf("[CheckRoleExist] role信息已存在")
		return errors.New("role信息已存在")
	}
	return nil
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
func AddRolePyDep(pyDep PyDep) (int64, error) {
	result := DB.Model(&PyDep{}).Create(&pyDep)
	if result.Error != nil {
		log.Printf("[AddRolePyDep] 数据库创建角色pyDep信息失败")
		return 0, result.Error
	}
	return pyDep.Id, nil
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
func GetRoleCode(roleId int64) (Code, error) {
	var codeInfo Code
	result := DB.Model(&Code{}).Where("role_id = ?", roleId).First(&codeInfo)
	if result.Error != nil {
		log.Printf("[GetRoleCode] 数据库获取角色code信息失败")
		return codeInfo, result.Error
	}
	log.Printf("[GetRoleCode] 数据库获取角色code信息成功")
	return codeInfo, nil
}

// GetRolePyDep 获取角色pyDep信息
func GetRolePyDep(roleId int64) (PyDep, error) {
	var pyDepInfo PyDep
	result := DB.Model(&PyDep{}).Where("role_id = ?", roleId).First(&pyDepInfo)
	if result.Error != nil {
		log.Printf("[GetRolePyDep] 数据库获取角色pyDep信息失败")
		return pyDepInfo, result.Error
	}
	log.Printf("[GetRolePyDep] 数据库获取角色pyDep信息成功")
	return pyDepInfo, nil
}

// GetRoleImage 获取角色image信息
func GetRoleImage(roleId int64) (Image, error) {
	var imageInfo Image
	result := DB.Model(&Image{}).Where("role_id = ?", roleId).First(&imageInfo)
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

// UpdateRoleInfo 更新角色信息
func UpdateRoleInfo(roleId int64, description, pyVersion, workDir, runCommand string) error {
	result := DB.Model(&Role{}).Where("id = ? ", roleId).
		Updates(Role{Description: description, PyVersion: pyVersion, WorkDir: workDir, RunCommand: runCommand})
	if result.Error != nil {
		log.Printf("[UpdateRoleInfo] 数据库更新角色信息失败")
		return result.Error
	}
	return nil
}

// UpdateRoleCode 更新角色code信息
func UpdateRoleCode(roleId, codeFileSize int64, codeSource, codeFileName, codeFileUrl, codeGitUrl string) error {
	if codeSource == "upload" {
		if codeFileName == "" {
			return nil
		}
		result := DB.Model(&Code{}).Where("role_id = ?", roleId).
			Updates(Code{CodeFileName: codeFileName, CodeFileSize: codeFileSize, CodeFileUrl: codeFileUrl})
		if result.Error != nil {
			log.Printf("[UpdateRoleCode] 数据库更新角色code信息失败")
			return result.Error
		}
		return nil
	} else if codeSource == "git" {
		result := DB.Model(&Code{}).Where("role_id = ?", roleId).
			Updates(Code{CodeGitUrl: codeGitUrl})
		if result.Error != nil {
			log.Printf("[UpdateRoleCode] 数据库更新角色code信息失败")
			return result.Error
		}
		return nil
	}
	return nil
}

// UpdateRolePyDep 更新角色pyDep信息
func UpdateRolePyDep(roleId int64, pyDepSource, pyDepPackage, pyDepGitUrl, pyDepGitFilepath string) error {
	if pyDepSource == "upload" || pyDepSource == "manual" {
		result := DB.Model(&PyDep{}).Where("role_id = ?", roleId).
			Updates(PyDep{PyDepPackages: pyDepPackage})
		if result.Error != nil {
			log.Printf("[UpdateRolePyDep] 数据库更新角色pyDep信息失败")
			return result.Error
		}
		return nil
	} else if pyDepSource == "git" {
		result := DB.Model(&PyDep{}).Where("role_id = ?", roleId).
			Updates(PyDep{PyDepGitUrl: pyDepGitUrl, PyDepGitFilepath: pyDepGitFilepath})
		if result.Error != nil {
			log.Printf("[UpdateRolePyDep] 数据库更新角色pyDep信息失败")
			return result.Error
		}
		return nil
	}
	return nil
}

// UpdateRoleImage 更新角色image信息
//TODO：imageName的更新未确定
func UpdateRoleImage(roleId int64, imageSource, imageName, imageDockerfileUrl, imageDockerfileName,
	imageArchiveUrl, imageArchiveName, imageGitUrl, imageGitFilepath string, imageDockerfileSize, imageArchiveSize int64) error {
	if imageSource == "platform" || imageSource == "dockerHub" {
		result := DB.Model(&Image{}).Where("role_id = ?", roleId).
			Updates(Image{ImageName: imageName})
		if result.Error != nil {
			log.Printf("[UpdateRoleImage] 数据库更新角色image_platform/dockerHub信息失败")
			return result.Error
		}
		result = DB.Model(&Role{}).Where("id = ? ", roleId).
			Updates(Role{ImageName: imageName})
		if result.Error != nil {
			log.Printf("[UpdateRoleImage] 数据库更新角色imageName信息失败")
			return result.Error
		}
		return nil
	} else if imageSource == "uploadDockerfile" {
		result := DB.Model(&Image{}).Where("role_id = ?", roleId).
			Updates(Image{ImageDockerfileUrl: imageDockerfileUrl, ImageDockerfileSize: imageDockerfileSize,
				ImageDockerfileName: imageDockerfileName})
		if result.Error != nil {
			log.Printf("[UpdateRoleImage] 数据库更新角色image_uploadDockerfile信息失败")
			return result.Error
		}
		return nil
	} else if imageSource == "uploadArchive" {
		result := DB.Model(&Image{}).Where("role_id = ?", roleId).
			Updates(Image{ImageArchiveUrl: imageArchiveUrl, ImageArchiveSize: imageArchiveSize,
				ImageArchiveName: imageArchiveName})
		if result.Error != nil {
			log.Printf("[UpdateRoleImage] 数据库更新角色image_uploadArchive信息失败")
			return result.Error
		}
		return nil
	} else if imageSource == "git" {
		result := DB.Model(&Image{}).Where("role_id = ?", roleId).
			Updates(Image{ImageGitUrl: imageGitUrl, ImageGitFilepath: imageGitFilepath})
		if result.Error != nil {
			log.Printf("[UpdateRoleImage] 数据库更新角色image_git信息失败")
			return result.Error
		}
		return nil
	}
	return nil
}

// DeleteRoleInfo 删除角色信息
func DeleteRoleInfo(roleId int64) error {
	result := DB.Where("id = ?", roleId).Delete(&Role{})
	if result.Error != nil {
		log.Printf("[DeleteRoleInfo] 数据库删除角色信息失败")
		return result.Error
	}
	return nil
}

// DeleteRoleCode 删除角色code信息
func DeleteRoleCode(roleId int64) error {
	result := DB.Where("role_id = ?", roleId).Delete(&Code{})
	if result.Error != nil {
		log.Printf("[DeleteRoleCode] 数据库删除角色code信息失败")
		return result.Error
	}
	return nil
}

// DeleteRolePyDep 删除角色pyDep信息
func DeleteRolePyDep(roleId int64) error {
	result := DB.Where("role_id = ?", roleId).Delete(&PyDep{})
	if result.Error != nil {
		log.Printf("[DeleteRolePyDep] 数据库删除角色pyDep信息失败")
		return result.Error
	}
	return nil
}

// DeleteRoleImage 删除角色image信息
func DeleteRoleImage(roleId int64) error {
	result := DB.Where("role_id = ?", roleId).Delete(&Image{})
	if result.Error != nil {
		log.Printf("[DeleteRoleImage] 数据库删除角色image信息失败")
		return result.Error
	}
	return nil
}

// DeleteRoleOutputItem 删除角色outputItem信息
func DeleteRoleOutputItem(roleId int64) error {
	result := DB.Where("role_id = ?", roleId).Delete(&OutputItem{})
	if result.Error != nil {
		log.Printf("[DeleteRoleOutputItem] 数据库删除角色outputItem信息失败")
		return result.Error
	}
	return nil
}
