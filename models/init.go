package models

import (
	"app/library/log"
)

func Init() {
	initDB()
	initCfg()
}

func initCfg() {
	if err := ReloadCfg(); err != nil {
		log.StdFatal("init", "initCfg", err)
	}
}
