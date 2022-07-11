package models

import (
	"app/library/log"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
	initDB()
	initCfg()
}

func initCfg() {
	if err := ReloadCfg(); err != nil {
		log.StdFatal("init", "initCfg", err)
	}
}

func SystemRunTime() int64 {
	return time.Now().Unix() - startTime.Unix()
}
