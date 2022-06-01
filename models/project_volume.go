package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
)

//k8s cluster virtual project ci/cd file
type ProjectVolume struct {
	Type    string `json:"type"` //url or content
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t *ProjectVolume) Validator() error {
	if t.Path == "" {
		return errors.New("file path is nil")
	}
	if t.Content == "" {
		return errors.New("file content is nil")
	}
	return nil
}

//k8s project cfg about file
type ProjectVolumeList []ProjectVolume

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *ProjectVolumeList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t ProjectVolumeList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *ProjectVolume) Load() (res string, err error) {
	if t.Type == "url" {
		resp, err := req.Get(t.Content)
		if err != nil {
			return
		}
		res = resp.String()
	} else if t.Type == "content" {
		res = t.Content
	}
	return
}
