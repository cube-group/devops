package models

import (
	"app/library/core"
	"app/library/log"
	"app/library/types/jsonutil"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os/exec"
	"time"
)

var _cfg CfgStruct

type CfgStruct struct {
	RegistryHost      string `json:"registryHost"`
	RegistryUsername  string `json:"registryUsername"`
	RegistryPassword  string `json:"registryPassword"`
	RegistryNamespace string `json:"registryNamespace"`
	Ci                []Kv   `json:"ci"`
	OnlineBlock       string `json:"onlineBlock"`
}

func (t *CfgStruct) Map() (res map[string]interface{}) {
	json.Unmarshal(jsonutil.ToBytes(t), &res)
	return
}

func (t *CfgStruct) Validator() error {
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
	//restart docker pull ci
	go func() {
		time.Sleep(10 * time.Second)
		for _, item := range _cfg.Ci {
			log.StdOut("docker pull start: " + item.V)
			log.StdOut("docker pull end: "+item.V, exec.Command("sh", "-c", "docker pull "+item.V).Run())
		}
	}()
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

func CfgGetCiList() []Kv {
	return _cfg.Ci
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
