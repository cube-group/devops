package log

import (
	"app/library/crypt/md5"
	"app/library/types/times"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

//日志启用时间戳
var _logStartTime = time.Now().Unix()
//当前正在操作日志的指针实例
var _logger *log.Logger
//当前正在写入日志的文件指针
var _currentFile *os.File
//唯一id
var _ruid = ""
//应用名称
var AppName string = "default"
//日志目录
var LogPath string = "/data/log/"

//初始化文件类日志系统
func initLog() {
	now := time.Now()
	newName := fmt.Sprintf("%s/%s-golib-%d.txt", path.Dir(LogPath), now.Format("2006-01-02"), _logStartTime)
	if _ruid == "" {
		_ruid = md5.MD5(fmt.Sprintf("%v-%v", rand.Intn(10000), time.Now().Nanosecond()))
	}

	if _currentFile != nil {
		if _currentFile.Name() != newName {
			_currentFile.Close()
		} else {
			return
		}
	}

	f, err := os.Create(newName)
	if err != nil {
		fmt.Println(err)
		_logger = log.New(os.Stdout, "", 0)
	} else {
		_currentFile = f
		_logger = log.New(f, "", 0)
	}
}

//记录本地文本类Info日志
func FileLogInfo(route, uid, code, msg, ext interface{}) {
	func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("[Log] Err", e)
			}
		}()

		initLog()
		_logger.Println(fmt.Sprintf(
			"%s|%s|%s|%s|%s|%v|%v|%v|%v|%v|%v",
			times.FormatDatetime(time.Now()),
			"127.0.0.1",
			"INFO",
			_ruid,
			AppName,
			0,
			route,
			uid,
			code,
			msg,
			ext,
		))
	}()
}

//记录本地文件类错误日志
func FileLogErr(route, uid, code, msg, ext interface{}) {
	func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("[Log] Err", e)
			}
		}()

		initLog()
		_logger.Println(fmt.Sprintf(
			"%s|%s|%s|%s|%s|%v|%v|%v|%v|%v|%v",
			times.FormatDatetime(time.Now()),
			"127.0.0.1",
			"ERROR",
			_ruid,
			AppName,
			0,
			route,
			uid,
			code,
			msg,
			ext,
		))
	}()
}
