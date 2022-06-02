package models

import (
	"app/library/core"
	"app/library/types/convert"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Cfg struct {
	ID    uint32 `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name  string `gorm:"" json:"name" form:"name" binding:"required"`
	Value string `gorm:"" json:"value" form:"value" binding:"required"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Cfg) TableName() string {
	return "d_cfg"
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

func CfgRegistryHost()string{
	return GetCfgString("registryHost")
}

func CfgRegistryUsername()string{
	return GetCfgString("registryUsername")
}

func CfgRegistryPassword()string{
	return GetCfgString("registryPassword")
}

func CfgRegistryNamespace()string{
	return GetCfgString("registryNamespace")
}

