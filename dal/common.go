package dal

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	UserName  string    `json:"user_name,omitempty" gorm:"type:varchar(100)"`
	Pwd       string    `json:"password,omitempty" gorm:"type:varchar(100)"`
	Salt      string    `json:"salt,omitempty" gorm:"type:char(4)"`
	Email     string    `json:"email,omitempty" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
	Dataset   []Dataset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role      []Role    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Dataset struct {
	Id          int64     `json:"id,omitempty" gorm:"primaryKey"`
	DatasetName string    `json:"dataset_name,omitempty" gorm:"type:varchar(100)"`
	Description string    `json:"description,omitempty" gorm:"type:varchar(100)"`
	Type        string    `json:"type,omitempty" gorm:"type:varchar(100)"`
	State       int8      `json:"state,omitempty" gorm:"type:int"`
	Source      string    `json:"source,omitempty" gorm:"type:varchar(100)"`
	Url         string    `json:"url,omitempty" gorm:"type:varchar(100)"`
	FileName    string    `json:"file_name,omitempty" gorm:"type:varchar(100)"`
	Size        int64     `json:"size,omitempty" gorm:"type:int"`
	UserId      int64     `json:"user_id,omitempty" gorm:"type:foreignKey"`
	CreatedAt   time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Role struct {
	Id          int64     `json:"id,omitempty" gorm:"primaryKey"`
	RoleName    string    `json:"role_name,omitempty" gorm:"type:varchar(100)"`
	Description string    `json:"description,omitempty" gorm:"type:varchar(100)"`
	PyVersion   string    `json:"py_version,omitempty" gorm:"type:varchar(100)"`
	WorkDir     string    `json:"work_dir,omitempty" gorm:"type:varchar(100)"`
	RunCommand  string    `json:"run_command,omitempty" gorm:"type:varchar(100)"`
	CodeId      int64     `json:"code_id,omitempty" gorm:"type:int"`
	PyDevId     int64     `json:"py_dev_id,omitempty" gorm:"type:int"`
	ImageId     int64     `json:"image_id,omitempty" gorm:"type:int"`
	UserId      int64     `json:"user_id,omitempty" gorm:"type:foreignKey"`
	CreatedAt   time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Code struct {
	Id           int64     `json:"id,omitempty" gorm:"primaryKey"`
	CodeSource   string    `json:"code_source,omitempty" gorm:"type:varchar(100)"`
	CodeFileUrl  string    `json:"code_file_url,omitempty" gorm:"type:varchar(100)"`
	CodeFileName string    `json:"code_file_name,omitempty" gorm:"type:varchar(100)"`
	CodeFileSize int64     `json:"code_file_size,omitempty" gorm:"type:int"`
	CodeGitUrl   string    `json:"code_git_url,omitempty" gorm:"type:varchar(100)"`
	CreatedAt    time.Time `json:"-" gorm:"index:,sort:desc"`
}

type PyDev struct {
	Id               int64     `json:"id,omitempty" gorm:"primaryKey"`
	PyDevSource      string    `json:"py_dev_source,omitempty" gorm:"type:varchar(100)"`
	PyDevPackages    string    `json:"py_dev_packages,omitempty" gorm:"type:varchar(1000)"`
	PyDevGitUrl      string    `json:"py_dev_git_url,omitempty" gorm:"type:varchar(100)"`
	PyDevGitFilepath string    `json:"py_dev_git_filepath,omitempty" gorm:"type:varchar(100)"`
	CreatedAt        time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Image struct {
	Id                  int64     `json:"id,omitempty" gorm:"primaryKey"`
	ImageSource         string    `json:"image_source,omitempty" gorm:"type:varchar(100)"`
	ImageName           string    `json:"image_name,omitempty" gorm:"type:varchar(100)"`
	ImageDockerfileUrl  string    `json:"image_dockerfile_url,omitempty" gorm:"type:varchar(100)"`
	ImageDockerfileName string    `json:"image_dockerfile_name,omitempty" gorm:"type:varchar(100)"`
	ImageDockerfileSize int64     `json:"image_dockerfile_size,omitempty" gorm:"type:int"`
	ImageArchiveUrl     string    `json:"image_archive_url,omitempty" gorm:"type:varchar(100)"`
	ImageArchiveName    string    `json:"image_archive_name,omitempty" gorm:"type:varchar(100)"`
	ImageArchiveSize    int64     `json:"image_archive_size,omitempty" gorm:"type:int"`
	ImageGitUrl         string    `json:"image_git_url,omitempty" gorm:"type:varchar(100)"`
	ImageGitFilepath    string    `json:"image_git_filepath,omitempty" gorm:"type:varchar(100)"`
	CreatedAt           time.Time `json:"-" gorm:"index:,sort:desc"`
}

type PlatformImage struct {
	Id          int64     `json:"id,omitempty" gorm:"primaryKey"`
	ImageName   string    `json:"image_name,omitempty" gorm:"type:varchar(100)"`
	Description string    `json:"description,omitempty" gorm:"type:varchar(100)"`
	CreatedAt   time.Time `json:"-" gorm:"index:,sort:desc"`
}

type OutputItem struct {
	Id         int64     `json:"id,omitempty" gorm:"primaryKey"`
	OutputName string    `json:"output_name,omitempty" gorm:"type:varchar(100)"`
	OutputPath string    `json:"output_path,omitempty" gorm:"type:varchar(100)"`
	RoleId     int64     `json:"role_id,omitempty" gorm:"type:int"`
	CreatedAt  time.Time `json:"-" gorm:"index:,sort:desc"`
}
