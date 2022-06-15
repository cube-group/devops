package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var callerPath string

//获取当前应用程序执行路径
//go build所在临时目录
func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//获取当前文件所在目录路径
func GetCallerPath(files ...string) (result string) {
	if callerPath == "" {
		_, fileStr, _, _ := runtime.Caller(1)
		callerPath = filepath.Dir(fileStr)
	}
	return GetFilePath(callerPath, files...)
}

func GetFilePath(dir string, files ...string) (result string) {
	result = dir
	if len(files) > 0 {
		for _, v := range files {
			result = filepath.Join(result, v)
		}
	}
	return
}
