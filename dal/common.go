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
	Project   []Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	ImageName   string    `json:"image_name,omitempty" gorm:"type:varchar(100)"`
	ProjectId   int64     `json:"project_id,omitempty" gorm:"type:int"`
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
	RoleId       int64     `json:"role_id,omitempty" gorm:"type:int"`
	CreatedAt    time.Time `json:"-" gorm:"index:,sort:desc"`
}

type PyDep struct {
	Id               int64     `json:"id,omitempty" gorm:"primaryKey"`
	PyDepSource      string    `json:"py_dep_source,omitempty" gorm:"type:varchar(100)"`
	PyDepPackages    string    `json:"py_dep_packages,omitempty" gorm:"type:varchar(1000)"`
	PyDepGitUrl      string    `json:"py_dep_git_url,omitempty" gorm:"type:varchar(100)"`
	PyDepGitFilepath string    `json:"py_dep_git_filepath,omitempty" gorm:"type:varchar(100)"`
	RoleId           int64     `json:"role_id,omitempty" gorm:"type:int"`
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
	RoleId              int64     `json:"role_id,omitempty" gorm:"type:int"`
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

type Project struct {
	Id                    int64     `json:"id,omitempty" gorm:"primaryKey"`
	ProjectName           string    `json:"project_name,omitempty" gorm:"type:varchar(100)"`
	DatasetId             int64     `json:"dataset_id,omitempty" gorm:"type:int"`
	ManagerFileId         int64     `json:"manager_file_id,omitempty" gorm:"type:int"`
	StructureFileId       int64     `json:"structure_file_id,omitempty" gorm:"type:int"`
	DatasetSplitterFileId int64     `json:"dataset_splitter_file_id,omitempty" gorm:"type:int"`
	UserId                int64     `json:"user_id,omitempty" gorm:"type:foreignKey"`
	CreatedAt             time.Time `json:"-" gorm:"index:,sort:desc"`
}

type File struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	Url       string    `json:"url,omitempty" gorm:"type:varchar(100)"`
	Name      string    `json:"name,omitempty" gorm:"type:varchar(100)"`
	Size      int64     `json:"size,omitempty" gorm:"type:int"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Config struct {
	Id           int64     `json:"id,omitempty" gorm:"primaryKey"`
	LinkType     string    `json:"link_type,omitempty" gorm:"type:varchar(100)"`
	BandwidthMax int64     `json:"bandwidth_max,omitempty" gorm:"type:int"`
	BandwidthMin int64     `json:"bandwidth_min,omitempty" gorm:"type:int"`
	ProjectId    int64     `json:"project_id,omitempty" gorm:"type:int"`
	CreatedAt    time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Node struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	NodeName  string    `json:"link_type,omitempty" gorm:"type:varchar(100)"`
	CPU       int64     `json:"cpu,omitempty" gorm:"type:int"`
	RAM       int64     `json:"ram,omitempty" gorm:"type:int"`
	RoleName  int64     `json:"role_name,omitempty" gorm:"type:varchar(100)"`
	ConfigId  int64     `json:"config_id,omitempty" gorm:"type:int"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Log struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	NodeName  string    `json:"node_name,omitempty" gorm:"type:varchar(100)"`
	Content   int64     `json:"content,omitempty" gorm:"type:varchar(1000)"`
	ProjectId int64     `json:"project_id,omitempty" gorm:"type:int"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}
