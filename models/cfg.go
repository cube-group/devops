package models

import (
	"app/library/core"
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
	ImageSsh          string `json:"imageSsh"`
	Ci                []Kv   `json:"ci"`
	OnlineBlock       string `json:"onlineBlock"`
}

func (t *CfgStruct) Map() (res map[string]string) {
	json.Unmarshal(jsonutil.ToBytes(t), &res)
	return
}

func (t *CfgStruct) Validator() error {
	if t.ImageSsh == "" {
		return errors.New("image ssh不能为空")
	}
	for i := 0; i < len(t.Ci); {
		if t.Ci[i].Validator() != nil {
			t.Ci = append(t.Ci[:i], t.Ci[i+1:]...)
		} else {
			i++
		}
	}
	if len(t.Ci) == 0 {
		return errors.New("ci镜像不能为空，至少保证有一个可用")
	}
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
	json.Unmarshal(jsonutil.ToBytes(kv), &res)
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
	//ssh image
	if err = createCfg("imageSsh", "cubegroup/devops-ssh"); err != nil {
		return
	}
	//ci default image
	res, err := GetCfgKey("ci")
	if err != nil {
		return
	}
	var list []Kv
	if json.Unmarshal([]byte(res), &list) == nil {
		if len(list) > 0 {
			return
		}
	}
	list = []Kv{{Key: "default", Value: "cubegroup/devops-ci-java"}}
	return createCfg("ci", jsonutil.ToString(list))
}
