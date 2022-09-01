package json

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type OrmJsonObject map[string]interface{}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *OrmJsonObject) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(b, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t OrmJsonObject) Value() (driver.Value, error) {
	return json.Marshal(t)
}


type OrmJsonList []interface{}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *OrmJsonList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(b, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t OrmJsonList) Value() (driver.Value, error) {
	return json.Marshal(t)
}
