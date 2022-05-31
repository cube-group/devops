package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ProjectTemplateMarshalJSON ProjectTemplate

//k8s project cfg about spec template
type ProjectTemplate struct {
	HealthCheck string `json:"healthCheck"`
	Cronjob     string `json:"cronjob"`
}

//override marshal json
func (t ProjectTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProjectTemplateMarshalJSON
	}{
		ProjectTemplateMarshalJSON(t),
	})
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *ProjectTemplate) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t ProjectTemplate) Value() (driver.Value, error) {
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
