package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ProjectTemplateDockerMarshalJSON ProjectTemplateDocker

//k8s project cfg about spec template
type ProjectTemplateDocker struct {
	Shell      string     `json:"shell"`
	Image      string     `json:"image"`
	Dockerfile string     `json:"dockerfile"`
	RunOptions string     `json:"runOptions"`
	Volume     VolumeList `gorm:"" json:"volume" form:"volume"`
}

func (t *ProjectTemplateDocker) Validator() error {
	for i := 0; i < len(t.Volume); {
		if t.Volume[i].Validator() != nil {
			t.Volume = append(t.Volume[:i], t.Volume[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}

//override marshal json
func (t ProjectTemplateDocker) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProjectTemplateDockerMarshalJSON
	}{
		ProjectTemplateDockerMarshalJSON(t),
	})
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *ProjectTemplateDocker) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t ProjectTemplateDocker) Value() (driver.Value, error) {
	return json.Marshal(t)
}
