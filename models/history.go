package models

import (
	"app/library/consts"
	"app/library/core"
	"app/library/crypt/md5"
	"app/library/log"
	"app/library/types/times"
	"context"
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
	"sync"
	"time"
)

type HistoryStatus string

const (
	HistoryStatusDefault HistoryStatus = ""
	HistoryStatusFailed  HistoryStatus = "error"
	HistoryStatusSuccess HistoryStatus = "ok"
	HistoryStatusDeleted HistoryStatus = "deleted"
)

//上线缓存
var onlineMaps sync.Map

func GetHistory(values ...interface{}) (res *History) {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			if vv <= 0 {
				return
			}
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

type HistoryMarshalJSON History

type History struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"-"`
	Uid       *uint32        `gorm:"" json:"uid" form:"-"`
	NodeId    uint32         `gorm:"" json:"nodeId" form:"nodeId"`
	Node      *Node          `gorm:"" json:"node" form:"node" binding:"-"`
	Version   string         `gorm:"" json:"version" form:"version"`
	Desc      string         `gorm:"" json:"desc" form:"desc" binding:"required"`
	Status    HistoryStatus  `gorm:"status" json:"status" form:"-"`
	ProjectId uint32         `gorm:"" json:"projectId" form:"-"`
	Project   *Project       `gorm:"" json:"project" form:"-"`
	Log       string         `gorm:"" json:"-" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	onlineCtx       context.Context
	onlineCtxCancel context.CancelFunc
}

func (t *History) TableName() string {
	return "d_history"
}

func (t History) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		HistoryMarshalJSON
	}{
		HistoryMarshalJSON(t),
	})
}

func (t *History) Validator() error {
	if t.Version == "" {
		t.Version = "latest"
	}
	if t.ProjectId > 0 {
		if i := GetProject(t.ProjectId); i != nil {
			t.Project = i
		} else {
			return fmt.Errorf("部署失败：project id %d not found", t.ProjectId)
		}
	}
	if t.Project.Mode != ProjectModeImage {
		if t.NodeId == 0 {
			return errors.New("部署目标机器不能为空")
		}
		t.Version = "latest" //强制改为latest
		if i := GetNode(t.NodeId); i != nil {
			t.Node = i
		} else {
			return fmt.Errorf("部署失败：node id %d not found", t.NodeId)
		}
	}
	return nil
}

func (t *History) ImageURL() string {
	return fmt.Sprintf(
		"%s/%s/%s:%s",
		CfgRegistryHost(), CfgRegistryNamespace(), t.Project.Name, t.Version,
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

func (t *History) Shutdown() error {
	if t.Status != HistoryStatusDefault {
		return errors.New("shutdown failed online is finished")
	}
	i, ok := onlineMaps.Load(t.ID)
	if ok {
		var target = i.(*History)
		if target.onlineCtxCancel != nil {
			target.onlineCtxCancel()
		} else {
			return errors.New("online no cancel")
		}
	} else {
		t.updateStatus(nil)
	}
	return nil
}

//project online
func (t *History) Online() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	t.onlineCtx = ctx
	t.onlineCtxCancel = cancel
	//store online history
	onlineMaps.Store(t.ID, t)

	var node = t.Node
	if t.Project.Mode != ProjectModeImage && node == nil {
		return errors.New("node not found")
	}
	//create workspace
	var workspace = t.Workspace()
	if err = os.MkdirAll(workspace, os.ModePerm); err != nil {
		return
	}
	//create run-dev.sh content
	var runContent string
	if t.Project.Mode == ProjectModeDocker || t.Project.Mode == ProjectModeImage { //deploy mode docker
		runContent, err = t.createRunDockerMode(node)
	} else { //node shell
		runContent, err = t.createRunNativeMode(node)
	}
	if err != nil {
		return
	}
	//create run-dev.sh
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
		cmd := exec.Command("sh", "-e", t.WorkspaceRun())
		cmd.Stdout = fileStream
		cmd.Stderr = fileStream
		t.updateStatus(cmd)
		cancel()
	}()
	return nil
}

func (t *History) updateStatus(cmd *exec.Cmd) {
	var err error
	var shutdownLogs string
	if cmd != nil {
		err = cmd.Start()
		if err == nil {
			go func() {
				select {
				case <-t.onlineCtx.Done():
					shutdownLogs = "shutdown 1"
					log.StdOut("history", "online", "SIGKILL", t.WorkspaceRun(), core.KillProcessGroup(cmd))
				}
			}()
			err = cmd.Wait()
		}
	} else { //shutdown
		shutdownLogs = "shutdown 2"
		err = errors.New("shutdown online")
	}

	var status = HistoryStatusSuccess
	if err != nil {
		status = HistoryStatusFailed
	}
	logBytes, _ := ioutil.ReadFile(t.WorkspaceFollowLog())
	defer func() {
		if er := os.RemoveAll(t.Workspace()); er != nil {
			log.StdWarning("history", "online", "clear workspace", t.Workspace(), er)
		}
	}()
	maps := map[string]interface{}{
		"status": status,
		"log":    string(logBytes) + "\r\n" + shutdownLogs,
	}
	if er := DB().Model(t).Where("id=?", t.ID).Updates(maps).Error; er != nil {
		log.StdWarning("history", "updateStatus", er)
	}
	//delete online history tag
	onlineMaps.Delete(t.ID)
}

func (t *History) createRunDockerMode(node *Node) (runContent string, err error) {
	var template = t.Project.Docker
	//create volumeLines
	var volumeLines = make([]string, 0)
	if t.Project.Docker.Image == "" {
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
	}

	var imageName string
	var dockerBuild string
	var dockerfile = template.Dockerfile
	if t.Project.Docker.Image != "" {
		imageName = t.Project.Docker.Image
	} else if dockerfile != "" { //create newLines
		var newLines = make([]string, 0)
		var fromImage string
		for _, v := range strings.Split(dockerfile, "\n") {
			v = strings.TrimLeft(v, " ")
			v = strings.TrimRight(v, " ")
			if strings.Contains(v, "FROM ") {
				fromImage = strings.Split(v, "FROM ")[1]
				v = fmt.Sprintf("%s\n%s", v, strings.Join(volumeLines, "\n"))
			}
			newLines = append(newLines, v)
		}
		if fromImage == "" {
			err = errors.New("Dockerfile invalid")
			return
		}
		dockerfile = strings.Join(newLines, "\n")
		//create dockerfile
		if err = ioutil.WriteFile(t.WorkspaceDockerfile(), []byte(dockerfile), os.ModePerm); err != nil {
			return
		}
		imageName = t.ImageURL()
		dockerBuild = fmt.Sprintf(`
cd %s
docker login %s --username=%s --password=%s
docker pull %s
docker build --platform=linux/amd64 -t %s . 
docker push %s
`,
			t.Workspace(),
			CfgRegistryHost(), CfgRegistryUsername(), CfgRegistryPassword(),
			fromImage,
			imageName, imageName)
	} else {
		err = errors.New("dockerfile or image is nil")
		return
	}

	//docker run
	var sshDockerRun string
	if t.Project.Mode == ProjectModeDocker {
		dockerRun := fmt.Sprintf(
			"docker login %s --username=%s --password=%s;"+
				"docker pull %s;"+
				"docker rm -f %s >/dev/null 2>&1;"+
				"docker run -it -d --restart=always --name %s %s %s",
			CfgRegistryHost(), CfgRegistryUsername(), CfgRegistryPassword(),
			imageName,
			t.Project.Name,
			t.Project.Name, t.Project.Docker.RunOptions, imageName,
		)
		sshDockerRun = fmt.Sprintf(
			"sshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'",
			node.SshPassword, node.SshPort, node.SshUsername, node.IP,
			dockerRun,
		)
	}

	//create run.sh content
	runContent = fmt.Sprintf(`
#!/bin/sh
date +"%%Y-%%m-%%d %%H:%%M:%%S"
cd %s
%s
%s
%s
cd %s
`,
		t.Workspace(),
		t.Project.Docker.Shell,
		dockerBuild,
		sshDockerRun,
		t.Workspace(),
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
	var shell = fmt.Sprintf("#!/bin/sh\n%s\n", template.Shell)
	if err = ioutil.WriteFile(shellFilePath, []byte(shell), os.ModePerm); err != nil {
		return
	}
	scpLines = append(scpLines, fmt.Sprintf(
		`sshpass -p '%s' scp -P %s -o StrictHostKeyChecking=no %s %s@%s:%s`,
		node.SshPassword, node.SshPort, shellFilePath,
		node.SshUsername, node.IP, tmpFilePath,
	))
	//create run-dev.sh content
	runContent = fmt.Sprintf(`
#!/bin/sh
echo 'ProjectId: %d HistoryId: %d %s'
echo '%s'
cd %s
%s
sshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s "sh %s"
`,
		t.Project.ID, t.ID, times.FormatDatetime(now), tmpFilePath,
		t.Workspace(),
		strings.Join(scpLines, "\n"),
		node.SshPassword, node.SshPort, node.SshUsername, node.IP,
		tmpFilePath,
	)
	return
}

func (t *History) Remove() error {
	if t.Project.Mode != ProjectModeDocker {
		return nil
	}
	_, err := t.Node.Exec(fmt.Sprintf("docker rm -f %s", t.Project.Name))
	return err
}
