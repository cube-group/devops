package models

import (
	"app/library/api"
	"app/library/core"
	"app/library/log"
	"app/library/types/jsonutil"
	"encoding/json"
	"errors"
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
	GitlabAddress     string `json:"gitlabAddress"`
	GitlabAppId       string `json:"gitlabAppId"`
	GitlabAppSecret   string `json:"gitlabAppSecret"`
	GitlabRedirectUri string `json:"gitlabRedirectUri"`
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

func Gitlab() (res *api.Gitlab, err error) {
	if _cfg.GitlabAddress == "" || _cfg.GitlabAppId == "" || _cfg.GitlabAppSecret == "" || _cfg.GitlabRedirectUri == "" {
		err = errors.New("gitlab option is not enough")
		return
	}
	res = api.NewGitlab(api.GitlabOption{
		GitlabAddress:     _cfg.GitlabAddress,
		GitlabAppId:       _cfg.GitlabAppId,
		GitlabAppSecret:   _cfg.GitlabAppSecret,
		GitlabRedirectUri: _cfg.GitlabRedirectUri,
	})
	return
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
