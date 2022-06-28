package models

import (
	"app/library/consts"
	"app/library/sshtool"
	"app/library/types/convert"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"strings"
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
	t.IP = strings.Trim(t.IP, " ")
	t.SshPort = strings.Trim(t.SshPort, " ")
	t.SshPassword = strings.Trim(t.SshPassword, " ")
	return nil
}

func (t *Node) WorkspacePath(ele ...string) string {
	ele = append([]string{"/devops/workspace/.node", convert.MustString(t.ID)}, ele...)
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

func (t *Node) initReadyIdRsa() (sshPath, idRsaPath string, err error) {
	if t.SshKey != "" {
		sshPath = t.WorkspacePath(".ssh")
		if err = os.MkdirAll(sshPath, 0700); err != nil {
			return
		}
		idRsaPath = t.WorkspacePath(".ssh/id_rsa")
		if err = ioutil.WriteFile(idRsaPath, []byte(t.SshKey), 0600); err != nil {
			return
		}
	}
	return
}

//node ssh args
func (t *Node) RunSshArgs(idRsaPath, remoteShell string) (args []string, err error) {
	if idRsaPath == "" {
		_, idRsaPath, err = t.initReadyIdRsa()
		if err != nil {
			return
		}
	}
	if idRsaPath != "" {
		args = []string{"ssh", "-i", idRsaPath}
	} else {
		args = []string{"sshpass", "-P", fmt.Sprintf("'%s'", t.SshPassword), "ssh"}
	}
	args = append(args, []string{
		"-p",
		t.SshPort,
		"-o",
		"StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", t.SshUsername, t.IP),
	}...)
	if remoteShell != "" {
		args = append(args, fmt.Sprintf("'%s'", remoteShell))
	}
	return
}

//node scp args
func (t *Node) RunScpArgs(localPath, remotePath string) (args []string, err error) {
	_, idRsaPath, err := t.initReadyIdRsa()
	if err != nil {
		return
	}
	if idRsaPath != "" {
		args = []string{"scp", "-i", idRsaPath}
	} else {
		args = []string{"sshpass", "-P", fmt.Sprintf("'%s'", t.SshPassword), "scp"}
	}
	args = append(args, []string{
		"-P",
		t.SshPort,
		"-o",
		"StrictHostKeyChecking=no",
		localPath,
		fmt.Sprintf("%s@%s:%s", t.SshUsername, t.IP, remotePath),
	}...)
	return
}
