package models

import (
	"app/library/core"
	"app/library/log"
	"app/library/types/convert"
	"fmt"
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

func ReloadCfg() error {
	res, err := GetCfg()
	if err != nil {
		return err
	}
	core.Lock(func() {
		cfgList = res
	})
	return nil
}

func GetCfg() (res gin.H, err error) {
	res = make(gin.H)
	var list []Cfg
	if err = DB().Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		res[v.Name] = v.Value
	}
	fmt.Println(res)
	return
}

func GetCfgString(key string) string {
	if v, ok := cfgList[key]; ok {
		return convert.MustString(v)
	}
	return ""
}
