package models

import (
	"app/library/consts"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func GetProject(values ...interface{}) *Project {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i Project
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextProject); exist {
				if instance, ok := i.(*Project); ok {
					return instance
				}
			}
			return nil
		}
	}
	return nil
}

type ProjectMarshalJSON Project

//virtual project
type Project struct {
	ID        uint32                `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name      string                `gorm:"index;column:name" json:"name" form:"name" binding:"required"`
	Desc      string                `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	Uid       *uint32               `gorm:"column:uid" json:"uid" form:"-"`
	Ding      string                `gorm:"column:ding" json:"ding" form:"ding"`
	Shell     string                `gorm:"column:shell" json:"shell" form:"shell"`
	Docker    ProjectDockerTemplate `gorm:"" json:"docker" form:"docker"`
	Volume    ProjectVolumeList     `gorm:"" json:"volume" form:"volume"`
	Cronjob   string                `gorm:"" json:"cronjob" form:"cronjob"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt gorm.DeletedAt        `gorm:"index" json:"-"`
}

func (t *Project) TableName() string {
	return "d_project"
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Project) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t Project) Value() (driver.Value, error) {
	return json.Marshal(t)
}

//override marshal json
func (t Project) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProjectMarshalJSON
	}{
		ProjectMarshalJSON(t),
	})
}

func (t *Project) GetLatestHistory() *History {
	var i History
	if DB().Last(&i, "project_id=?", t.ID).Error != nil {
		return nil
	}
	return &i
}

//k8s cluster visual project data validator check
func (t *Project) Validator() error {
	for i := 0; i < len(t.Volume); {
		if t.Volume[i].Validator() != nil {
			t.Volume = append(t.Volume[:i], t.Volume[i+1:]...)
		} else {
			i++
		}
	}
	if err := t.Docker.Validator(); err != nil {
		return err
	}
	return nil
}
