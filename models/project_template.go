package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ProjectDockerTemplateMarshalJSON ProjectDockerTemplate

//k8s project cfg about spec template
type ProjectDockerTemplate struct {
	Dockerfile   string      `json:"dockerfile"`
	Name         string      `json:"name"`
	Privileged   uint32      `json:"privileged"`
	Cmd          string      `json:"cmd"`
	Ports        []ProjectKv `json:"ports"`
	Volumes      []ProjectKv `json:"volumes"`
	Environments []ProjectKv `json:"environments"`
}

func(t *ProjectDockerTemplate)Validator()error{
	for i := 0; i < len(t.Ports); {
		if t.Ports[i].Validator() != nil {
			t.Ports = append(t.Ports[:i], t.Ports[i+1:]...)
		} else {
			i++
		}
	}
	for i := 0; i < len(t.Volumes); {
		if t.Volumes[i].Validator() != nil {
			t.Volumes = append(t.Volumes[:i], t.Volumes[i+1:]...)
		} else {
			i++
		}
	}
	for i := 0; i < len(t.Environments); {
		if t.Environments[i].Validator() != nil {
			t.Environments = append(t.Environments[:i], t.Environments[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}

//override marshal json
func (t ProjectDockerTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProjectDockerTemplateMarshalJSON
	}{
		ProjectDockerTemplateMarshalJSON(t),
	})
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *ProjectDockerTemplate) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t ProjectDockerTemplate) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type ProjectKv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (t *ProjectKv) Validator() error {
	if t.Key == "" || t.Value == "" {
		return errors.New("key or value is nil")
	}
	return nil
}
