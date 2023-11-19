package handler

type FileReturn struct {
	Url      string `form:"url" json:"url"`           //文件url
	FileName string `form:"fileName" json:"fileName"` //文件名称
	Size     int    `form:"size" json:"size"`         //文件大小
}

// CheckZipFile 检查文件类型是否符合
func CheckZipFile(uploadFileType string) bool {
	fileTypeList := []string{".zip", ".rar", ".gz", ".tar.gz", "tgz", "bz2", "z", "tar"}
	for _, element := range fileTypeList {
		if uploadFileType == element {
			return true
		}
	}
	return false
}
