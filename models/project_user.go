package models

import (
	"gorm.io/gorm"
	"time"
)

type ProjectUser struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Uid       uint32         `gorm:"" json:"uid" form:"-"`
	Pid       uint32         `gorm:"" json:"pid" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *ProjectUser) TableName() string {
	return "d_project_user"
}

func GetUserProject(uid uint32) (res []Project) {
	DB().Find(&res, "uid=?", uid)
	return
}
