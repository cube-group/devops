package models

import (
	"app/library/consts"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type ProjectMode string

type ProjectAccessLevel string

const (
	ProjectModeNative ProjectMode = "native"
	ProjectModeDocker ProjectMode = "docker"
	ProjectModeImage  ProjectMode = "image"

	ProjectAccessLevelNone   = "none"
	ProjectAccessLevelMember = "member"
)

func GetProject(values ...interface{}) (res *Project) {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i Project
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			res = &i
		case string:
			var i Project
			if err := DB().Take(&i, "name LIKE ?", "%"+vv+"%").Error; err != nil {
				return nil
			}
			res = &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextProject); exist {
				if instance, ok := i.(*Project); ok {
					res = instance
				}
			}
		}
	}
	//rel tag
	if res != nil {
		var rel TagRel
		if DB().Last(&rel, "pid=?", res.ID).Error == nil {
			if tag := GetTag(rel.Tid); tag != nil {
				res.Tag = tag.ID
			}
		}
	}
	return
}

type ProjectList []Project

func (t ProjectList) IDs() []uint32 {
	var list = make([]uint32, 0)
	for _, v := range t {
		list = append(list, v.ID)
	}
	return list
}

type ProjectMarshalJSON Project

//virtual project
type Project struct {
	ID        uint32                `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name      string                `gorm:"index;column:name" json:"name" form:"name" binding:"required"`
	Desc      string                `gorm:"" json:"desc" form:"desc" binding:"required"`
	Uid       uint32                `gorm:"" json:"uid" form:"-"`
	Ding      string                `gorm:"" json:"ding" form:"ding"`
	Mode      ProjectMode           `gorm:"" json:"mode" form:"mode"`
	Native    ProjectTemplateNative `gorm:"" json:"native" form:"native"`
	Docker    ProjectTemplateDocker `gorm:"" json:"docker" form:"docker"`
	Cronjob   string                `gorm:"" json:"cronjob" form:"cronjob"`
	Deleted   uint32                `gorm:"" json:"deleted" form:"-"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt gorm.DeletedAt        `gorm:"index" json:"-"`

	Tag uint32 `gorm:"-" json:"tag" form:"tag"`
}

func (t *Project) TableName() string {
	return "d_project"
}

// ?????? sql.Scanner ?????????Scan ??? value ????????? Jsonb
func (t *Project) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// ?????? driver.Valuer ?????????Value ?????? json value
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

func (t *Project) GetLatestHistory(option ...interface{}) *History {
	var i History
	if DB(option...).Last(&i, "project_id=?", t.ID).Error != nil {
		return nil
	}
	return &i
}

//k8s cluster visual project data validator check
func (t *Project) Validator() error {
	if matched, err := regexp.MatchString("^[a-z0-9-]{4,40}$", t.Name); err != nil || !matched {
		return errors.New("?????????????????????????????????^[a-z0-9-]{4,40}$")
	}
	if err := t.Native.Validator(); err != nil {
		return err
	}
	if err := t.Docker.Validator(); err != nil {
		return err
	}
	if t.Mode == ProjectModeDocker {
		if t.Docker.Dockerfile == "" && t.Docker.Image == "" {
			return errors.New("docker????????????Dockerfile???Image??????????????????")
		}
	} else if t.Mode == ProjectModeNative {
		if t.Native.Shell == "" {
			return errors.New("native????????????Shell????????????")
		}
	}
	return nil
}

func (t *Project) IsCronjob() bool {
	return t.Cronjob != ""
}

func (t *Project) Apply(history *History, async bool) (err error) {
	if err = DB().Transaction(func(tx *gorm.DB) error {
		//Block online
		if tx.Last(&History{}, "project_id=? AND status=?", t.ID, HistoryStatusDefault).Error == nil {
			return errors.New("????????????????????????????????????????????????...")
		}
		//Update project
		if er := tx.Model(t).Where("id=?", t.ID).Update("deleted", 0).Error; er != nil {
			return er
		}
		//Save history
		if er := tx.Save(history).Error; er != nil {
			return er
		}
		return nil
	}); err != nil {
		return
	}
	return history.Online(async)
}

func (t *Project) StopCronjob() (err error) {
	return CronjobStop(t.ID)
}

//???????????????????????????
func (t *Project) RollbackVersions() (res []History) {
	var list HistoryList
	var db = DB()
	if db.Select("max(id) AS id").Order("id DESC").Group("version").Where("project_id=? AND status=?", t.ID, HistoryStatusSuccess).Find(&list).Error == nil {
		db.Order("id DESC").Find(&res, "id IN (?)", list.IDs())
	}
	return
}
