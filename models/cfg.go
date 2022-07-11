package models

import (
	"app/library/core"
	"app/library/log"
	"app/library/types/jsonutil"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

var _cfg CfgStruct

type CfgStruct struct {
	RegistryHost      string `json:"registryHost"`
	RegistryUsername  string `json:"registryUsername"`
	RegistryPassword  string `json:"registryPassword"`
	RegistryNamespace string `json:"registryNamespace"`
	OnlineBlock       string `json:"onlineBlock"`
}

func (t *CfgStruct) Map() (res map[string]interface{}) {
	json.Unmarshal(jsonutil.ToBytes(t), &res)
	return
}

func (t *CfgStruct) Validator() error {
	return nil
}

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
	i, err := GetCfg()
	if err != nil {
		log.StdWarning("cfg", "reload", err)
		return err
	}
	core.Lock(func() {
		_cfg = i
	})
	return nil
}

func GetCfg() (res CfgStruct, err error) {
	kv := make(gin.H)
	var list []Cfg
	if err = DB().Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		var i interface{}
		if json.Unmarshal([]byte(v.Value), &i) != nil {
			kv[v.Name] = v.Value
		} else {
			kv[v.Name] = i
		}
	}
	err = json.Unmarshal(jsonutil.ToBytes(kv), &res)
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

func CfgRegistryPath() (res string) {
	s, err := GetCfg()
	if err != nil {
		return
	}
	return fmt.Sprintf("%s/%s", s.RegistryHost, s.RegistryNamespace)
}

func CfgOnlineBlock() string {
	res, _ := GetCfgKey("onlineBlock")
	return res
}

func createCfg(key, value string) (err error) {
	var find Cfg
	if DB().Take(&find, "`name`=?", key).Error != nil {
		find = Cfg{Name: key, Value: value}
	} else {
		find.Value = value
	}
	return DB().Save(&find).Error
}

func createDefaultCfg() (err error) {
	//ci default image
	var list []Kv
	res, err := GetCfgKey("ci")
	if err == nil {
		if json.Unmarshal([]byte(res), &list) == nil {
			if len(list) > 0 {
				return
			}
		}
	}

	list = []Kv{
		{K: "default", V: "cubegroup/devops-ci-default"},
		{K: "java", V: "cubegroup/devops-ci-java"},
	}
	if err = createCfg("ci", jsonutil.ToString(list)); err != nil {
		return
	}
	return nil
}
