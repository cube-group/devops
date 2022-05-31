package models

import (
	"gorm.io/gorm"
	"time"
)

type Cfg struct {
	ID    uint32 `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name  string `gorm:"" json:"name" form:"name" binding:"required"`
	Value string `gorm:"" json:"value" form:"value" binding:"required"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Cfg) TableName() string {
	return "d_cfg"
}
