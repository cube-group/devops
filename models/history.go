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

func GetHistory(values ...interface{}) *History {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i History
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextHistory); exist {
				if instance, ok := i.(*History); ok {
					return instance
				}
			}
			return nil
		}
	}
	return nil
}

func GetHistoryProject(c *gin.Context) *Project {
	project := GetProject(c)
	if project == nil {
		return project
	}
	var projectLatestHistory History
	if DB().Last(&projectLatestHistory, "pid=?", project.ID).Error == nil {
		return projectLatestHistory.Project
	}
	return project
}

type HistoryTime struct {
	StartTime   int64 `json:"startTime"`
	CiStartTime int64 `json:"ciStartTime"`
	CiStopTime  int64 `json:"ciStopTime"`
	CdStartTime int64 `json:"cdStartTime"`
	CdStopTime  int64 `json:"cdStopTime"`
	StopTime    int64 `json:"stopTime"`
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *HistoryTime) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t HistoryTime) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type HistoryMarshalJSON History

type History struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Uid       *uint32        `gorm:"" json:"uid" form:"-"`
	Pid       uint32         `gorm:"" json:"pid" form:"-"`
	Desc      string         `gorm:"" json:"desc" form:"desc" binding:"required"`
	Time      HistoryTime   `gorm:"time" json:"time" form:"-"`
	Project   *Project      `gorm:"" json:"project" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *History) TableName() string {
	return "c_history"
}
func (t History) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		HistoryMarshalJSON
	}{
		HistoryMarshalJSON(t),
	})
}

func (t *History) Validator() error {
	return nil
}

//save new history
//project apply to k8s
//then watch k8s apply pipeline
func (t *History) Save() error {
	return DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(t).Error; err != nil {
			return err
		}
		//go t.Watch() // watch k8s apply pipeline
		return nil
	})
}