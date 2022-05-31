// Author: chenqionghe
// Time: 2018-10
// 自定义错误

package ginutil

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

//自定义错误
type MyError struct {
	code int
	msg  string
}

type MyForceError struct {
	Err *MyError
}

//新建错误对象
func Error(msg ...interface{}) *MyError {
	var array []string
	var msgStr string
	for _, v := range msg {
		//如果是validator的错误，定义相应格式
		if errs, ok := v.(validator.ValidationErrors); ok {
			msgStr = formatValidateError(errs)
		} else {
			msgStr = fmt.Sprint(v)
		}
		array = append(array, msgStr)
	}
	return &MyError{
		code: CODE_ERR,
		msg:  strings.Join(array, ""),
	}
}

//新建强提示错误
func ForceError(msg ...interface{}) *MyForceError {
	err := Error(msg...)
	return &MyForceError{&MyError{
		code: CODE_FORCE_ERR,
		msg:  err.Error(),
	}}
}
func (t *MyForceError) String() string {
	return t.Err.msg
}
func (t *MyForceError) Error() string {
	return t.Err.String()
}

func (t *MyError) String() string {
	return t.msg
}

func (t *MyError) Error() string {
	return t.String()
}

func (t *MyError) Code() int {
	return t.code
}

//返回validator错误信息友好格式
func formatValidateError(errs validator.ValidationErrors) string {
	var array []string
	fieldErrMsg := "%s不能通过%s规则"
	for _, v := range errs {
		//array = append(array, fmt.Sprintf(fieldErrMsg, str.Camel2Case(v.Field()), v.Tag()))
		array = append(array, fmt.Sprintf(fieldErrMsg, v.Field(), v.Tag()))
	}
	msgStr := strings.Join(array, "<br>")
	return msgStr
}
