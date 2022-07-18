package models

import (
	"app/library/consts"
	"app/library/log"
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

func GetSomeNodes(ids []uint32) (res NodeList, err error) {
	var db = DB()
	err = db.Find(&res, "id IN (?)", ids).Error
	if err == nil {
		for _, id := range ids {
			if !res.Has(id) {
				err = fmt.Errorf("node id: %d not found", id)
				return
			}
		}
	}
	return
}

func GetNodes() (res []Node) {
	DB().Order("id DESC").Find(&res)
	return res
}

func NodeClean() {
	for _, item := range GetNodes() {
		_, err := item.Exec("docker system prune -a -f")
		log.StdOut("nodeClean", item.IP, err)
	}
}

type NodeList []Node

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *NodeList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t NodeList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t NodeList) Has(id uint32) bool {
	for _, v := range t {
		if v.ID == id {
			return true
		}
	}
	return false
}

func (t NodeList) Get(id uint32) (node Node, ok bool) {
	if id == 0 {
		return
	}
	for _, v := range t {
		if v.ID == id {
			return v, true
		}
	}
	return
}

func (t NodeList) Del(id uint32) NodeList {
	for k, v := range t {
		if v.ID == id {
			return append(t[:k], t[k+1:]...)
		}
	}
	return t
}

func (t NodeList) Contact(node Node) NodeList {
	for _, v := range t {
		if v.ID == node.ID {
			return t
		}
	}
	return append(t, node)
}

type NodeMarshalJSON Node

//virtual Node
type Node struct {
	ID          uint32 `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name        string `gorm:"" json:"name" form:"name" binding:"required"`
	Desc        string `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	IP          string `gorm:"" json:"ip" form:"ip" binding:"required"` //内网IP
	Uid         uint32 `gorm:"" json:"uid" form:"-"`
	SshPort     string `gorm:"" json:"sshPort" form:"required"`
	SshKey      string `gorm:"" json:"sshKey" form:"required"` //private key
	SshUsername string `gorm:"" json:"sshUsername" form:"required"`
	SshPassword string `gorm:"" json:"sshPassword" form:"required"`

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

func (t *Node) IsNone() error {
	if t.ID == 0 {
		return errors.New("node is nil")
	}
	return nil
}

func (t *Node) WorkspacePath(ele ...string) string {
	ele = append([]string{"/devops/workspace/.node", convert.MustString(t.ID)}, ele...)
	return path.Join(ele...)
}

func (t *Node) WorkspaceSshPath() string {
	return t.WorkspacePath(".ssh")
}

func (t *Node) WorkspaceSshIdRsaPath() string {
	return t.WorkspacePath(".ssh", "id_rsa")
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

func (t *Node) initWorkspace() (err error) {
	if t.SshKey != "" {
		if err = os.MkdirAll(t.WorkspaceSshPath(), 0700); err != nil {
			return
		}
		if err = ioutil.WriteFile(t.WorkspaceSshIdRsaPath(), []byte(t.SshKey), 0600); err != nil {
			return
		}
	}
	return
}

//node ssh args
func (t *Node) RunSshArgs(tty bool, idRsaPath, remoteShell string) (args []string, err error) {
	if err = t.IsNone(); err != nil {
		return
	}
	if idRsaPath == "" {
		err = t.initWorkspace()
		if err != nil {
			return
		}
		idRsaPath = t.WorkspaceSshIdRsaPath()
	}
	if t.SshKey != "" {
		args = []string{"ssh", "-i", idRsaPath}
	} else {
		args = []string{"sshpass", "-p", fmt.Sprintf("%s", t.SshPassword), "ssh"}
	}
	if tty {
		args = append(args, "-t")
	}
	args = append(args, []string{
		"-p",
		t.SshPort,
		"-o",
		"StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", t.SshUsername, t.IP),
	}...)
	if remoteShell != "" {
		args = append(args, fmt.Sprintf("%s", remoteShell))
	}
	return
}

//node scp args
//container special
func (t *Node) RunScpArgs(localPath, remotePath string) (args []string, err error) {
	if err = t.IsNone(); err != nil {
		return
	}
	err = t.initWorkspace()
	if err != nil {
		return
	}
	if t.SshKey != "" {
		args = []string{"scp", "-i", t.WorkspaceSshIdRsaPath()} //containerSsh2Path
	} else {
		args = []string{"sshpass", "-p", fmt.Sprintf("'%s'", t.SshPassword), "scp"}
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
