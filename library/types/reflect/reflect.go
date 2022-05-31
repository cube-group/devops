package reflect

import (
	"app/library/e"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//struct转换为map
func GetStructMap(target interface{}) (map[string]interface{}, error) {
	var err error
	var data = make(map[string]interface{})
	func() {
		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()
		t := reflect.TypeOf(target)
		v := reflect.ValueOf(target)
		for i := 0; i < t.NumField(); i++ {
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}()

	return data, err
}

// 获取struct中的所有字段名称
func GetStructFields(target interface{}) ([]string, error) {
	var res []string
	var err error
	func() {
		defer func() {
			if e := recover(); e != nil { //catch e
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()

		t := reflect.TypeOf(target)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			panic("Check type e not Struct")
		}

		fieldNum := t.NumField()
		res = make([]string, 0, fieldNum)
		for i := 0; i < fieldNum; i++ {
			res = append(res, t.Field(i).Name)
		}
	}()

	return res, err
}

// 获取struct中的所有字段名称
func GetStructTags(target interface{}) ([]string, error) {
	var res []string
	var err error
	func() {
		defer func() {
			if e := recover(); e != nil { //catch e
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()

		t := reflect.TypeOf(target)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			panic("Check type e not Struct")
		}
		fieldNum := t.NumField()
		res = make([]string, 0, fieldNum)
		for i := 0; i < fieldNum; i++ {
			tagName := t.Field(i).Name
			tags := strings.Split(string(t.Field(i).Tag), "\"")
			if len(tags) > 1 {
				tagName = tags[1]
			}
			res = append(res, tagName)
		}
	}()

	return res, err
}

// 给struct里的字段自动赋值
func SetStructFieldValues(target interface{}, values map[string]interface{}) error {
	var err error
	func() {
		defer func() {
			if e := recover(); e != nil { //catch e
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()

		v := reflect.ValueOf(target).Elem() // the struct variable
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i) // a reflect.StructField
			tag := fieldInfo.Tag           // a reflect.StructTag
			name := tag.Get("json")
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
			name = strings.Split(name, ",")[0]

			if value, ok := values[name]; ok {
				//给结构体赋值
				//保证赋值时数据类型一致
				if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
					v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
				}
			}
		}
	}()

	return err
}

//返回v实例的类别
func TypeOfString(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func GetTagFromStruct(v interface{}, field, tag, defaultV string) (res string) {
	e.TryCatch(func() {
		s, _ := reflect.TypeOf(v).FieldByName(field)
		res = s.Tag.Get(tag)
	})
	if res == "" {
		res = defaultV
	}
	return
}

func GetTagFromStructPoint(v interface{}, field, tag, defaultV string) (res string) {
	e.TryCatch(func() {
		s, _ := reflect.TypeOf(v).Elem().FieldByName(field)
		res = s.Tag.Get(tag)
	})
	if res == "" {
		res = defaultV
	}
	return
}
