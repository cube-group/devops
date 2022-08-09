package models

import (
	"app/library/consts"
	"app/library/e"
	"app/library/log"
	"app/library/sshtool"
	"app/library/types/convert"
	"app/library/types/jsonutil"
	"app/library/types/number"
	"app/library/types/times"
	"bufio"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	//NodeContainerRandomPortStart Host container startup random start port
	//scope: [NodeContainerRandomPortStart, NodeContainerRandomPortStop)
	NodeContainerRandomPortStart = 51000
	NodeContainerRandomPortStop  = 60000
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

func (t NodeList) IsActive() bool {
	for _, v := range t {
		if !v.Removed {
			return true
		}
	}
	return false
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

func (t NodeList) RemovePod(id uint32) {
	for k, v := range t {
		if v.ID == id {
			t[k].Removed = true
		}
	}
}

func (t NodeList) Contact(node Node) NodeList {
	for _, v := range t {
		if v.ID == node.ID {
			return t
		}
	}
	return append(t, node)
}

func (t NodeList) Security() {
	for k, _ := range t {
		t[k].SshKey = ""
		t[k].SshPassword = ""
	}
}

//global node list proc info map
var nodeProcMap sync.Map

type NodeContainerPsItem map[string]interface{}

//0.0.0.0:8082->80/tcp, :::8082->80/tcp
func (t NodeContainerPsItem) Ports() (res map[int]bool) {
	res = make(map[int]bool)
	ports, ok := t["Ports"]
	if !ok {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(convert.MustString(ports)))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var text = scanner.Text()
		var publishArr = strings.Split(text, "->")
		if publishArr == nil || len(publishArr) < 2 {
			continue
		}
		var publishHostPortArr = strings.Split(publishArr[0], ":")
		if publishHostPort := convert.MustInt(publishHostPortArr[len(publishHostPortArr)-1]); publishHostPort > 0 {
			res[publishHostPort] = true
		}
	}
	return
}

type NodeMarshalJSON Node

//virtual Node
type Node struct {
	ID          uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name        string         `gorm:"" json:"name" form:"name" binding:"required"`
	Desc        string         `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	IP          string         `gorm:"" json:"ip" form:"ip" binding:"required"` //内网IP
	Uid         uint32         `gorm:"" json:"uid" form:"-"`
	SshPort     string         `gorm:"" json:"sshPort" form:"required"`
	SshKey      string         `gorm:"" json:"sshKey" form:"required"` //private key
	SshUsername string         `gorm:"" json:"sshUsername" form:"required"`
	SshPassword string         `gorm:"" json:"sshPassword" form:"required"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Removed     bool           `gorm:"-" json:"removed" form:"-" binding:"-"` //是否被移除
}

func (t *Node) TableName() string {
	return "d_node"
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Node) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(b, &t) //no error
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
		CreatedTime string      `json:"createdTime"`
		Proc        interface{} `json:"proc"`
	}{
		NodeMarshalJSON: NodeMarshalJSON(t),
		CreatedTime:     times.FormatDatetime(t.CreatedAt),
		Proc:            t.proc(),
	})
}

//k8s cluster visual Node data validator check
func (t *Node) Validator(c *gin.Context) error {
	t.IP = strings.Trim(t.IP, " ")
	t.SshPort = strings.Trim(t.SshPort, " ")
	t.SshPassword = strings.Trim(t.SshPassword, " ")
	return nil
}

//sync exec remote shell
func (t *Node) ExecMulti(cmdList []string) (res []sshtool.SSHClientSessionResult) {
	client, err := sshtool.NewSSHClient(t.IP, t.SshPort, t.SshUsername, t.SshPassword, t.SshKey)
	if err != nil {
		return
	}
	defer client.Close()
	return client.ExecMulti(cmdList)
}

//sync exec remote shell
func (t *Node) Exec(cmd string) (res []byte, err error) {
	client, err := sshtool.NewSSHClient(t.IP, t.SshPort, t.SshUsername, t.SshPassword, t.SshKey)
	if err != nil {
		return
	}
	defer client.Close()
	result := client.Exec(cmd)
	return result.Result, result.Error
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

//get node cache proc info
func (t *Node) proc() gin.H {
	if value, ok := nodeProcMap.Load(t.ID); ok {
		return value.(gin.H)
	} else {
		return gin.H{}
	}
}

//exec node proc process
func (t *Node) Proc() {
	var info gin.H
	if value, ok := nodeProcMap.Load(t.ID); ok {
		info = value.(gin.H)
	} else {
		info = gin.H{}
	}
	cmdList := []string{
		`cat /proc/cpuinfo | grep "processor"| wc -l`,
		`head -n 1 /proc/stat`,
		`cat /proc/meminfo`,
	}
	if e.TryCatch(func() {
		execList := t.ExecMulti(cmdList)
		if execCpuInfo := execList[0]; execCpuInfo.Error == nil {
			info["cpuNum"] = convert.MustInt32(strings.Trim(string(execCpuInfo.Result), "\n "))
		}
		if execCpuStat := execList[1]; execCpuStat.Error == nil {
			cpuStat, cpuPercent := t.procStat(string(execCpuStat.Result), info["cpuStat"])
			info["cpuStat"] = cpuStat
			info["cpuPercent"] = cpuPercent
		}
		if execMemInfo := execList[2]; execMemInfo.Error == nil {
			total, free := t.procMemInfo(string(execMemInfo.Result))
			info["memTotal"] = total
			info["memFree"] = free
		}
		nodeProcMap.Store(t.ID, info)
	}) != nil {
		return
	}
}

//analyse node /proc/meminfo
//MemTotal:       15858772 kB
//MemFree:         1314852 kB
//MemAvailable:    6439260 kB
func (t *Node) procMemInfo(content string) (total, free int64) {
	var scanner = bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lineScanner := bufio.NewScanner(strings.NewReader(line))
		lineScanner.Split(bufio.ScanWords)
		var lineKey string
		if lineScanner.Scan() {
			lineKey = lineScanner.Text()
		}
		if lineScanner.Scan() {
			if strings.Contains(lineKey, "MemTotal") {
				total = convert.MustInt64(strings.Trim(lineScanner.Text(), " ")) * 1024
			} else if strings.Contains(lineKey, "MemFree") {
				free = convert.MustInt64(strings.Trim(lineScanner.Text(), " ")) * 1024
			}
		}
		if total > 0 && free > 0 {
			break
		}
	}
	return
}

//analyse node /proc/stat
func (t *Node) procStat(content string, latestStatInter interface{}) (stat []float64, percent float64) {
	s := bufio.NewScanner(strings.NewReader(content))
	s.Split(bufio.ScanWords)
	if s.Scan() && s.Text() != "cpu" {
		return
	}
	stat = make([]float64, 0)
	for s.Scan() {
		stat = append(stat, convert.MustFloat64(s.Text()))
	}
	if latestStatInter != nil {
		var err = e.TryCatch(func() {
			var latestStat []float64
			if jsonutil.ToJson(latestStatInter, &latestStat) != nil {
				return
			}
			var statA = ((stat[0] + stat[1] + stat[2]) - (latestStat[0] + latestStat[1] + latestStat[2]))
			var statB = ((stat[0] + stat[1] + stat[2] + stat[3]) - (latestStat[0] + latestStat[1] + latestStat[2] + latestStat[3]))
			percent = number.NumberDecimal(statA/statB, 6)
		})
		if err != nil {
			log.StdWarning("node", "proc.stat", err)
		}
	}
	return
}

func (t *Node) GetDockerVersion() (bytes []byte, err error) {
	return t.Exec("docker version --format='{{json .}}'")
}

func (t *Node) GetDockerContainerList() (res []NodeContainerPsItem, err error) {
	bytes, err := t.Exec("docker ps --format='{{json .}}'")
	if err != nil {
		return
	}
	res = make([]NodeContainerPsItem, 0)
	for _, strItem := range strings.Split(string(bytes), "\n") {
		var resItem NodeContainerPsItem
		if jsoniter.Unmarshal([]byte(strItem), &resItem) == nil {
			res = append(res, resItem)
		}
	}
	return
}

//GetContainerRandomPort
//Get the random start port of the host container
func (t *Node) GetContainerRandomPort(length int) (ports []int, err error) {
	containerList, err := t.GetDockerContainerList()
	if err != nil {
		return
	}
	var hasPublishedPorts = make(map[int]bool)
	for _, item := range containerList {
		for port, _ := range item.Ports() {
			hasPublishedPorts[port] = true
		}
	}
	//循环分配新端口
	var newPublishPorts = make([]int, 0)
	for newPort := NodeContainerRandomPortStart; newPort < NodeContainerRandomPortStop; newPort++ {
		if _, ok := hasPublishedPorts[newPort]; !ok {
			newPublishPorts = append(newPublishPorts, newPort)
			if len(newPublishPorts) >= length {
				break
			}
		}
	}
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
