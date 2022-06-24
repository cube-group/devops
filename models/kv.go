package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Kv struct {
	K string `json:"key"`
	V string `json:"value"`
}

func (t *Kv) Validator() error {
	if t.K == "" || t.V == "" {
		return errors.New("key or value is nil")
	}
	return nil
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Kv) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t Kv) Value() (driver.Value, error) {
	return json.Marshal(t)
}
