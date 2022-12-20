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
}

type Dataset struct {
	Id          int64       `json:"id,omitempty" gorm:"primaryKey"`
	SetName     string      `json:"set_name,omitempty" gorm:"type:varchar(100)"`
	Type        string      `json:"type,omitempty" gorm:"type:varchar(100)"`
	Size        int64       `json:"size,omitempty" gorm:"type:int"`
	Description string      `json:"description,omitempty" gorm:"type:varchar(100)"`
	State       int8        `json:"state,omitempty" gorm:"type:int"`
	UserId      int64       `json:"user_id,omitempty" gorm:"type:foreignKey"`
	CreatedAt   time.Time   `json:"-" gorm:"index:,sort:desc"`
	DatasetFile DatasetFile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type DatasetFile struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	FileName  string    `json:"file_name,omitempty" gorm:"type:varchar(100)"`
	Path      string    `json:"path,omitempty" gorm:"type:varchar(100)"`
	Type      string    `json:"type,omitempty" gorm:"type:varchar(100)"`
	Size      int64     `json:"size,omitempty" gorm:"type:int"`
	DatasetId int64     `json:"dataset_id,omitempty" gorm:"type:foreignKey"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}
