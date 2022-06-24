package models

import (
	"app/library/consts"
	"app/library/log"
	"app/library/uuid"
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

	volumePath = ".volume"
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
	Ci        Kv             `gorm:"" json:"ci" form:"ci"`
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
	if !t.IsDockerImageMode() {
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
	if t.IsDockerMode() {
		if ci, ok := t.findCi(); !ok {
			return errors.New("构建器不能为空")
		} else {
			t.Ci = ci
		}
	}
	return nil
}

func (t *History) findCi() (res Kv, ok bool) {
	for _, item := range CfgGetCiList() {
		if item.K == t.Ci.K {
			return item, true
		}
	}
	return
}

func (t *History) IsDockerMode() bool {
	return t.Project.Mode == ProjectModeDocker || t.Project.Mode == ProjectModeImage
}

func (t *History) IsDockerImageMode() bool {
	return t.Project.Mode == ProjectModeImage
}

func (t *History) ImageURL() string {
	return fmt.Sprintf(
		"%s/%s/%s:%s",
		_cfg.RegistryHost, _cfg.RegistryNamespace, t.Project.Name, t.Version,
	)
}

func (t *History) Workspace() string {
	return fmt.Sprintf("/devops/workspace/%d/%d", t.Project.ID, t.ID)
}

func (t *History) WorkspacePath(elem ...string) string {
	elem = append([]string{t.Workspace()}, elem...)
	return path.Join(elem...)
}

func (t *History) WorkspaceNativePath(md5 string) string {
	return path.Join(t.Workspace(), volumePath, md5)
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
	//node check
	var node = t.Node
	if t.Project.Mode != ProjectModeImage && node == nil {
		return errors.New("node not found")
	}
	//create workspace
	var workspace = t.WorkspacePath(volumePath)
	if err = os.MkdirAll(workspace, os.ModePerm); err != nil {
		return
	}
	//ssh key ready
	var sshKeyVolume string
	if node != nil && node.SshKey != "" {
		if err = os.MkdirAll(t.WorkspacePath("/.ssh"), 0700); err != nil {
			return
		}
		if err = ioutil.WriteFile(t.WorkspacePath("/.ssh/id_rsa"), []byte(node.SshKey), 0600); err != nil {
			return
		}
		sshKeyVolume = fmt.Sprintf("-v %s:/root/.ssh", t.WorkspacePath(".ssh"))
	}
	//exec docker run --rm
	var execContent string
	if t.Project.Mode == ProjectModeDocker || t.Project.Mode == ProjectModeImage { //deploy mode docker
		execContent, err = t.createRunDockerExec(node, sshKeyVolume)
	} else { //node shell
		execContent, err = t.createRunNativeExec(node, sshKeyVolume)
	}
	if err != nil {
		return
	}
	//create .follow.log
	fileStream, err := os.Create(t.WorkspaceFollowLog())
	if err != nil {
		return
	}
	go func() {
		defer fileStream.Close()
		cmd := exec.Command("sh", "-c", execContent)
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
				}
			}()
			err = cmd.Wait()
		}
	} else { //shutdown
		shutdownLogs = "shutdown 2"
		err = errors.New("shutdown online")
	}
	t.dockerRunClose() //close docker run
	var status = HistoryStatusSuccess
	if err != nil {
		status = HistoryStatusFailed
	}
	logBytes, _ := ioutil.ReadFile(t.WorkspaceFollowLog())
	defer func() {
		time.Sleep(time.Second)
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

func (t *History) createRunDockerVolume(filePath string, content string) (copy string, err error) {
	var fileName = uuid.GetUUID(t.ID, "@", filePath)
	var localPath = t.WorkspaceNativePath(fileName)
	if err = ioutil.WriteFile(localPath, []byte(content), os.ModePerm); err != nil {
		return
	}
	copy = fmt.Sprintf("COPY %s %s", fileName, filePath)
	return
}

//${workspace}/Dockerfile -> ${docker}/Dockerfile
//${workspace}/.volume -> ${docker}/.volume
//${workspace}/run.sh -> ${docker}/run.sh
func (t *History) createRunDockerExec(node *Node, sshKeyVolume string) (execContent string, err error) {
	//TODO 需要检测之前的history如果存在project.name不一致需要先移除container
	var template = t.Project.Docker
	//create Dockerfile COPY
	var volumeLines = make([]string, 0)
	if t.Project.Docker.Image == "" {
		for _, v := range template.Volume {
			var volumeContent string
			volumeContent, err = v.Load()
			if err != nil {
				return
			}
			copyLine, er := t.createRunDockerVolume(v.Path, volumeContent)
			if er != nil {
				err = er
				return
			}
			volumeLines = append(volumeLines, copyLine)
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
		if err = ioutil.WriteFile(t.WorkspacePath(volumePath, "Dockerfile"), []byte(dockerfile), os.ModePerm); err != nil {
			return
		}
		imageName = t.ImageURL()
		dockerBuild = fmt.Sprintf(`
docker login %s --username=%s --password=%s
docker pull %s
docker build --platform=linux/amd64 -t %s /%s
docker push %s
`,
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			fromImage, imageName, volumePath,
			imageName)
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
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			imageName,
			t.Project.Name,
			t.Project.Name, t.Project.Docker.RunOptions, imageName,
		)
		if t.Node.SshKey != "" {
			sshDockerRun = fmt.Sprintf(
				"ssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'",
				node.SshPort, node.SshUsername, node.IP,
				dockerRun,
			)
		} else {
			sshDockerRun = fmt.Sprintf(
				"sshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s '%s'",
				node.SshPassword, node.SshPort, node.SshUsername, node.IP,
				dockerRun,
			)
		}
	}

	//${docker} /run.sh
	dockerRun := fmt.Sprintf(
		"#!/bin/sh\ncd /%s\n%s\n%s\n%s\n",
		volumePath,
		t.Project.Docker.Shell,
		dockerBuild,
		sshDockerRun,
	)
	if err = ioutil.WriteFile(t.WorkspacePath("run.sh"), []byte(dockerRun), os.ModePerm); err != nil {
		return
	}

	//exec
	execContent = fmt.Sprintf(
		"docker run -i --rm --name %s -v /var/run/docker.sock:/var/run/docker.sock -v %s:/%s -v %s:/run.sh %s %s",
		t.dockerRunRmName(),
		t.WorkspacePath(volumePath), volumePath,
		t.WorkspacePath("run.sh"),
		sshKeyVolume, "cubegroup/devops-ci-java",
	)
	return
}

func (t *History) createRunNativeVolume(runFlag bool, filePath string, content string) (scp string, err error) {
	var fileName string
	if runFlag {
		fileName = "run.sh"
	} else {
		fileName = uuid.GetUUID(t.ID, "@", filePath)
	}
	var localPath = t.WorkspaceNativePath(fileName)
	if err = ioutil.WriteFile(localPath, []byte(content), os.ModePerm); err != nil {
		return
	}
	var containerPath = path.Join("/"+volumePath, fileName)
	var remoteNodePath = fmt.Sprintf("/tmp/devops-%d-%d-%s", t.ProjectId, t.ID, fileName)
	if t.Node.SshKey != "" {
		scp = fmt.Sprintf(
			`scp -P %s -o StrictHostKeyChecking=no %s %s@%s:%s`,
			t.Node.SshPort, containerPath,
			t.Node.SshUsername, t.Node.IP, remoteNodePath,
		)
	} else {
		scp = fmt.Sprintf(
			`sshpass -p '%s' scp -P %s -o StrictHostKeyChecking=no %s %s@%s:%s`,
			t.Node.SshPassword, t.Node.SshPort, containerPath,
			t.Node.SshUsername, t.Node.IP, remoteNodePath,
		)
	}
	return
}

//${workspace}/.volume -> ${docker}/.volume
//${workspace}/run.sh -> ${docker}/run.sh
func (t *History) createRunNativeExec(node *Node, sshKeyVolume string) (execContent string, err error) {
	var template = t.Project.Native
	//create scp files
	var scpLines = make([]string, 0)
	for _, v := range template.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		scp, er := t.createRunNativeVolume(false, v.Path, volumeContent)
		if er != nil {
			err = er
			return
		}
		scpLines = append(scpLines, scp)
	}
	//create scp run.sh
	scp, err := t.createRunNativeVolume(true, "",
		fmt.Sprintf("#!/bin/sh\n%s\nrm -rf /tmp/devops-%d-%d-*", template.Shell, t.ProjectId, t.ID),
	)
	if err != nil {
		return
	}
	scpLines = append(scpLines, scp)
	//create docker run /run.sh
	var remoteNodePath = fmt.Sprintf("/tmp/devops-%d-%d-%s", t.ProjectId, t.ID, "run.sh")
	var remoteShContent string
	if t.Node.SshKey != "" {
		remoteShContent = fmt.Sprintf("#!/bin/sh\n%s\nssh -p %s -o StrictHostKeyChecking=no %s@%s 'sh -ex %s'",
			strings.Join(scpLines, "\n"), node.SshPort, node.SshUsername, node.IP, remoteNodePath,
		)
	} else {
		remoteShContent = fmt.Sprintf("#!/bin/sh\n%s\nsshpass -p '%s' ssh -p %s -o StrictHostKeyChecking=no %s@%s 'sh -ex %s'",
			strings.Join(scpLines, "\n"),
			node.SshPassword, node.SshPort, node.SshUsername, node.IP, remoteNodePath,
		)
	}
	if err = ioutil.WriteFile(t.WorkspacePath("run.sh"), []byte(remoteShContent), os.ModePerm); err != nil {
		return
	}

	//exec
	execContent = fmt.Sprintf(
		"docker run -i --rm --name %s -v %s:/%s -v %s:/run.sh %s %s",
		t.dockerRunRmName(),
		t.WorkspacePath(volumePath), volumePath,
		t.WorkspacePath("run.sh"),
		sshKeyVolume, t.Ci.V,
	)
	return
}

func (t *History) Remove() error {
	if t.Project.Mode != ProjectModeDocker {
		return nil
	}
	_, err := t.Node.Exec(fmt.Sprintf("docker rm -f %s", t.Project.Name))
	if err != nil {
		return err
	}
	return t.dockerRunClose()
}

func (t *History) dockerRunRmName() string {
	return fmt.Sprintf("devops-ci-%d-%d", t.ProjectId, t.ID)
}

func (t *History) dockerRunClose() error {
	return exec.Command("docker", "rm", "-f", t.dockerRunRmName()).Run()
}
