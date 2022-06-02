package models

import (
	"app/library/log"
	"github.com/gin-gonic/gin"
)

var cfgList gin.H

func Init() {
	initDB()
	initCfg()
}

func initCfg() {
	if err := ReloadCfg(); err != nil {
		log.StdFatal("init", "initCfg", err)
	}
}