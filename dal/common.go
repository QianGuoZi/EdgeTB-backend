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
