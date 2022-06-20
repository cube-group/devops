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
	return
}

func GetCfgKey(key string) (res string, err error) {
	var i Cfg
	if err = DB().Take(&i, "name=?", key).Error; err != nil {
		return
	}
	res = i.Value
	return
}

func GetCfgCache(key string) string {
	if v, ok := cfgList[key]; ok {
		return convert.MustString(v)
	}
	return ""
}

func CfgRegistryHost() string {
	return GetCfgCache("registryHost")
}

func CfgRegistryUsername() string {
	return GetCfgCache("registryUsername")
}

func CfgRegistryPassword() string {
	return GetCfgCache("registryPassword")
}

func CfgRegistryNamespace() string {
	return GetCfgCache("registryNamespace")
}

func CfgRegistryPath() string {
	return fmt.Sprintf("%s/%s", CfgRegistryHost(), CfgRegistryNamespace())
}

func CfgOnlineBlock() string {
	res, _ := GetCfgKey("onlineBlock")
	return res
}
