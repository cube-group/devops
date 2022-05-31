package env

import (
	"app/library/types/convert"
	"encoding/json"
	"os"
)

func Get(name string, defaultValue interface{}) interface{} {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func GetString(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func GetStringSlice(name string, defaultValue []string) []string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	var res []string
	if json.Unmarshal([]byte(v), &res) != nil {
		return nil
	}
	return res
}

func GetInt(name string, defaultValue int, scope ...int) int {
	var value int
	if v := os.Getenv(name); v != "" {
		value = convert.MustInt(v)
	} else {
		value = defaultValue
	}
	if len(scope) > 0 {
		if value < scope[0] {
			value = scope[0]
		}
	}
	if len(scope) > 1 {
		if value > scope[1] {
			value = scope[1]
		}
	}
	return value
}

func GetUint(name string, defaultValue uint, scope ...uint) uint {
	var value uint
	if v := os.Getenv(name); v != "" {
		value = convert.MustUint(v)
	} else {
		value = defaultValue
	}
	if len(scope) > 0 {
		if value < scope[0] {
			value = scope[0]
		}
	}
	if len(scope) > 1 {
		if value > scope[1] {
			value = scope[1]
		}
	}
	return value
}


func GetUint64(name string, defaultValue uint64, scope ...uint64) uint64 {
	var value uint64
	if v := os.Getenv(name); v != "" {
		value = convert.MustUint64(v)
	} else {
		value = defaultValue
	}
	if len(scope) > 0 {
		if value < scope[0] {
			value = scope[0]
		}
	}
	if len(scope) > 1 {
		if value > scope[1] {
			value = scope[1]
		}
	}
	return value
}

func GetInt64(name string, defaultValue int64, scope ...int64) int64 {
	var value int64
	if v := os.Getenv(name); v != "" {
		value = convert.MustInt64(v)
	} else {
		value = defaultValue
	}
	if len(scope) > 0 {
		if value < scope[0] {
			value = scope[0]
		}
	}
	if len(scope) > 1 {
		if value > scope[1] {
			value = scope[1]
		}
	}
	return value
}

func GetFloat64(name string, defaultValue float64, scope ...float64) float64 {
	var value float64
	if v := os.Getenv(name); v != "" {
		value = convert.MustFloat64(v)
	} else {
		value = defaultValue
	}
	if len(scope) > 0 {
		if value < scope[0] {
			value = scope[0]
		}
	}
	if len(scope) > 1 {
		if value > scope[1] {
			value = scope[1]
		}
	}
	return value
}



