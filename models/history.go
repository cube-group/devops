package models

import (
	"app/library/consts"
	"app/library/crypt/md5"
	"app/library/log"
	"app/library/types/times"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type HistoryStatus string

const (
	HistoryStatusDefault HistoryStatus = ""
	HistoryStatusFailed  HistoryStatus = "error"
	HistoryStatusSuccess HistoryStatus = "ok"
)

func GetHistory(values ...interface{}) (res *History) {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i History
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return
			}
			res = &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextHistory); exist {
				if instance, ok := i.(*History); !ok {
					return
				} else {
					res = instance
				}
			}
		}
	}
	return
}

type HistoryTime struct {
	StartTime   int64 `json:"startTime"`
	CiStartTime int64 `json:"ciStartTime"`
	CiStopTime  int64 `json:"ciStopTime"`
	CdStartTime int64 `json:"cdStartTime"`
	CdStopTime  int64 `json:"cdStopTime"`
	StopTime    int64 `json:"stopTime"`
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *HistoryTime) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t HistoryTime) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type HistoryMarshalJSON History

type History struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"-"`
	Uid       *uint32        `gorm:"" json:"uid" form:"-"`
	NodeId    uint32         `gorm:"" json:"nodeId" form:"nodeId" binding:"required"`
	Node      *Node          `gorm:"" json:"node" form:"node" binding:"-"`
	Version   string         `gorm:"" json:"version" form:"version" binding:"required"`
	Desc      string         `gorm:"" json:"desc" form:"desc" binding:"required"`
	Status    HistoryStatus  `gorm:"status" json:"status" form:"-"`
	ProjectId uint32         `gorm:"" json:"projectId" form:"-"`
	Project   *Project       `gorm:"" json:"project" form:"-"`
	Log       string         `gorm:"" json:"-" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *History) TableName() string {
	return "d_history"
}
func (t History) MarshalJSON() ([]byte, error) {
	if t.Project == nil {
		t.Project = GetProject(t.ProjectId)
	}
	if t.Node == nil {
		t.Node = GetNode(t.NodeId)
	}
	return json.Marshal(struct {
		HistoryMarshalJSON
	}{
		HistoryMarshalJSON(t),
	})
}

func (t *History) Validator() error {
	if t.ProjectId > 0 {
		t.Project = GetProject(t.ProjectId)
	}
	if t.NodeId > 0 {
		t.Node = GetNode(t.NodeId)
	}
	return nil
}

func (t *History) ImageURL() string {
	return fmt.Sprintf(
		"%s/%s/%s:latest",
		CfgRegistryHost(), CfgRegistryNamespace(), t.Project.Name,
	)
}

func (t *History) Workspace() string {
	return fmt.Sprintf("/devops/workspace/%d/%d", t.Project.ID, t.ID)
}

func (t *History) WorkspacePath(elem ...string) string {
	elem = append([]string{t.Workspace()}, elem...)
	return path.Join(elem...)
}

func (t *History) WorkspaceSshShellPath(md5 string) string {
	return path.Join(t.Workspace(), "ssh-file-"+md5)
}

func (t *History) WorkspaceFollowLog() string {
	return path.Join(t.Workspace(), ".follow.log")
}

func (t *History) WorkspaceEndLog() (res string) {
	if err := os.MkdirAll(t.Workspace(), os.ModePerm); err != nil {
		return
	}
	var endLogPath = path.Join(t.Workspace(), ".end.log")
	if err := ioutil.WriteFile(endLogPath, []byte(t.Log), os.ModePerm); err != nil {
		return
	}
	return endLogPath
}

func (t *History) WorkspaceRun() string {
	return path.Join(t.Workspace(), "run.sh")
}

func (t *History) WorkspaceDockerfile() string {
	return path.Join(t.Workspace(), "Dockerfile")
}

func (t *History) IsEnd() bool {
	return t.Status == HistoryStatusFailed || t.Status == HistoryStatusSuccess
}

//project online
func (t *History) Online() (err error) {
	var node = GetNode(t.Node.ID)
	if node == nil {
		return errors.New("node not found")
	}
	//create workspace
	var workspace = t.Workspace()
	if err = os.MkdirAll(workspace, os.ModePerm); err != nil {
		return
	}
	//create run.sh content
	var runContent string
	if t.Project.Mode == ProjectModeDocker { //deploy mode docker
		runContent, err = t.createRunDockerMode(node)
	} else { //node shell
		runContent, err = t.createRunNativeMode(node)
	}
	if err != nil {
		return
	}
	//create run.sh
	if err = ioutil.WriteFile(t.WorkspaceRun(), []byte(runContent), os.ModePerm); err != nil {
		return
	}
	//create .follow.log
	fileStream, err := os.Create(t.WorkspaceFollowLog())
	if err != nil {
		return
	}
	go func() {
		defer fileStream.Close()
		cmd := exec.Command("bash", t.WorkspaceRun())
		cmd.Stdout = fileStream
		cmd.Stderr = fileStream
		t.updateStatus(cmd.Run())
	}()
	return nil
}

func (t *History) updateStatus(err error) {
	var status = HistoryStatusSuccess
	if err != nil {
		status = HistoryStatusFailed
	}
	logBytes, _ := ioutil.ReadFile(t.WorkspaceFollowLog())
	defer func() {
		time.Sleep(5 * time.Minute)
		os.RemoveAll(t.Workspace())
	}()
	maps := map[string]interface{}{
		"status": status,
		"log":    string(logBytes),
	}
	if er := DB().Model(t).Where("id=?", t.ID).Updates(maps).Error; er != nil {
		log.StdWarning("history", "updateStatus", er)
	}
}

func (t *History) createRunDockerMode(node *Node) (runContent string, err error) {
	var template = t.Project.Docker
	var dockerfile = template.Dockerfile
	//create volumeLines
	var volumeLines = make([]string, 0)
	for k, v := range template.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		var volumeCopyFileName = md5.MD5(fmt.Sprintf("%d@%s", k, v.Path))
		if err = ioutil.WriteFile(t.WorkspacePath(volumeCopyFileName), []byte(volumeContent), os.ModePerm); err != nil {
			return
		}
		volumeLines = append(volumeLines, fmt.Sprintf("COPY %s %s", volumeCopyFileName, v.Path))
	}
	//create newLines
	var newLines = make([]string, 0)
	var fromFlag bool
	for _, v := range strings.Split(template.Dockerfile, "\n") {
		v = strings.TrimLeft(v, " ")
		v = strings.TrimRight(v, " ")
		if strings.Contains(v, "FROM ") {
			fromFlag = true
			v = fmt.Sprintf("%s\n%s", v, strings.Join(volumeLines, "\n"))
		}
		newLines = append(newLines, v)
	}
	if !fromFlag {
		err = errors.New("Dockerfile invalid")
		return
	}
	dockerfile = strings.Join(newLines, "\n")
	//create dockerfile
	if err = ioutil.WriteFile(t.WorkspaceDockerfile(), []byte(dockerfile), os.ModePerm); err != nil {
		return
	}
	var imageName = t.ImageURL()
	var dockerRun = fmt.Sprintf(
		"docker run -it -d --restart=always --name %s %s",
		t.Project.Name, imageName,
	)
	//create run.sh content
	runContent = fmt.Sprintf(`
#!/bin/bash
date +"%%Y-%%m-%%d %%H:%%M:%%S"
cd %s
docker login %s --username=%s --password=%s
docker build -t %s . 
docker push %s
sshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'
`,
		t.Workspace(),
		CfgRegistryHost(), CfgRegistryUsername(), CfgRegistryPassword(),
		imageName, imageName,
		node.SshPassword, node.SshPort, node.SshUsername, node.IP,
		dockerRun,
	)
	return
}

func (t *History) createRunNativeMode(node *Node) (runContent string, err error) {
	var now = time.Now()
	var template = t.Project.Native
	//create volumeLines
	var scpLines = make([]string, 0)
	for k, v := range template.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		var volumeFilePath = t.WorkspaceSshShellPath(md5.MD5(fmt.Sprintf("%d%s", k, v.Path)))
		if err = ioutil.WriteFile(volumeFilePath, []byte(volumeContent), os.ModePerm); err != nil {
			return
		}
		scpLines = append(scpLines, fmt.Sprintf(
			`sshpass -p '%s' scp -P %s -o StrictHostKeyChecking=no %s %s@%s:%s`,
			node.SshPassword, node.SshPort, volumeFilePath,
			node.SshUsername, node.IP, v.Path,
		))
	}
	//create ssh shell
	var shellFilePath = t.WorkspaceSshShellPath("")
	var tmpFilePath = fmt.Sprintf("/tmp/devops-%d-%s", t.Project.ID, times.FormatFileDatetime(now))
	var shell = fmt.Sprintf("#!/bin/bash\n%s\n", template.Shell)
	if err = ioutil.WriteFile(shellFilePath, []byte(shell), os.ModePerm); err != nil {
		return
	}
	scpLines = append(scpLines, fmt.Sprintf(
		`sshpass -p '%s' scp -P %s -o StrictHostKeyChecking=no %s %s@%s:%s`,
		node.SshPassword, node.SshPort, shellFilePath,
		node.SshUsername, node.IP, tmpFilePath,
	))
	//create run.sh content
	runContent = fmt.Sprintf(`
#!/bin/bash
echo 'ProjectId: %d HistoryId: %d %s'
echo '%s'
cd %s
%s
sshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s "bash %s"
`,
		t.Project.ID, t.ID, times.FormatDatetime(now), tmpFilePath,
		t.Workspace(),
		strings.Join(scpLines, "\n"),
		node.SshPassword, node.SshPort, node.SshUsername, node.IP,
		tmpFilePath,
	)
	return
}
