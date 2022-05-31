package log

import (
	"app/library/types/times"
	"fmt"
	"os"
	"time"
)

var enabled = true

func StdDisable(flag bool) {
	enabled = !flag
}

//打印标准输出流INFO日志
func StdOut(logType string, values ...interface{}) {
	if !enabled {
		return
	}
	values = append([]interface{}{"[" + logType + "]", times.FormatDatetime(time.Now()), "|", "INFO", "|"}, values...)
	fmt.Println(values...)
}

//打印标准输出流WARNING日志
func StdWarning(logType string, values ...interface{}) {
	if !enabled {
		return
	}
	values = append([]interface{}{"[" + logType + "]", times.FormatDatetime(time.Now()), "|", "WARN", "|"}, values...)
	fmt.Println(values...)
}

//打印标准输出流ERROR日志
func StdErr(logType string, values ...interface{}) {
	if !enabled {
		return
	}
	values = append([]interface{}{"[" + logType + "]", times.FormatDatetime(time.Now()), "|", "ERROR", "|"}, values...)
	fmt.Println(values...)
}

//打印标准输出流FATAL日志
func StdFatal(logType string, values ...interface{}) {
	if !enabled {
		return
	}
	values = append([]interface{}{"[" + logType + "]", times.FormatDatetime(time.Now()), "|", "FATAL", "|"}, values...)
	fmt.Println(values...)
	os.Exit(1)
}
