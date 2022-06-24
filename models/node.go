package models

import (
	"app/library/consts"
	"app/library/sshtool"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func GetNode(values ...interface{}) *Node {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i Node
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextNode); exist {
				if instance, ok := i.(*Node); ok {
					return instance
				}
			}
			return nil
		}
	}
	return nil
}

func GetNodes() (res []Node) {
	DB().Order("id DESC").Find(&res)
	return res
}

type NodeMarshalJSON Node

//virtual Node
type Node struct {
	ID          uint32  `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name        string  `gorm:"" json:"name" form:"name" binding:"required"`
	Desc        string  `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	IP          string  `gorm:"" json:"ip" form:"ip" binding:"required"` //内网IP
	Uid         *uint32 `gorm:"" json:"uid" form:"-"`
	SshPort     string  `gorm:"" json:"sshPort" form:"required"`
	SshKey      string  `gorm:"" json:"sshKey" form:"required"` //private key
	SshUsername string  `gorm:"" json:"sshUsername" form:"required"`
	SshPassword string  `gorm:"" json:"sshPassword" form:"required"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Node) TableName() string {
	return "d_node"
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Node) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t Node) Value() (driver.Value, error) {
	return json.Marshal(t)
}

//override marshal json
func (t Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeMarshalJSON
	}{
		NodeMarshalJSON(t),
	})
}

//k8s cluster visual Node data validator check
func (t *Node) Validator(c *gin.Context) error {
	return nil
}

func (t *Node) WorkspacePath(ele ...string) string {
	ele = append([]string{"/devops/workspace/node"}, ele...)
	return path.Join(ele...)
}

//sync exec remote shell
func (t *Node) Exec(cmd string) (res []byte, err error) {
	s, err := sshtool.SSHConnect(t.SshUsername, t.SshPassword, t.SshKey, t.IP, t.SshPort)
	if err != nil {
		return
	}
	defer s.Close()
	res, err = s.CombinedOutput(cmd)
	return
}

//node docker run
func (t *Node) DockerRun(name, mainCmd, otherCmd string, volumes map[string]string) (args []string, err error) {
	if volumes == nil {
		volumes = make(map[string]string)
	}
	var remoteShContent string
	if t.SshKey != "" {
		if err = os.MkdirAll(t.WorkspacePath(name, ".ssh"), 0700); err != nil {
			return
		}
		if err = ioutil.WriteFile(t.WorkspacePath(name, ".ssh/id_rsa"), []byte(t.SshKey), 0600); err != nil {
			return
		}
		volumes[t.WorkspacePath(name, ".ssh")] = "/root.ssh"
		remoteShContent = fmt.Sprintf("#!/bin/sh\n%s\nsshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'",
			otherCmd,
			t.SshPassword, t.SshPort, t.SshUsername, t.IP, mainCmd,
		)
	} else {
		remoteShContent = fmt.Sprintf("#!/bin/sh\n%s\nssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'",
			otherCmd, t.SshPort, t.SshUsername, t.IP, mainCmd,
		)
	}
	if err = ioutil.WriteFile(t.WorkspacePath(name, "run.sh"), []byte(remoteShContent), os.ModePerm); err != nil {
		return
	}
	volumes[t.WorkspacePath(name, "run.sh")] = "/run.ssh"
	args = []string{
		"docker",
		"run",
		"-it",
		"--rm",
		"--name",
		name,
	}
	for k, v := range volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", k, v))
	}
	args = append(args, _cfg.ImageSsh)
	return
}
