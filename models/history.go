package models

import (
	"app/library/consts"
	"app/library/core"
	"app/library/crypt/md5"
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

	DockerVolumePathName = ".volume"
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
	Uid       uint32         `gorm:"" json:"uid" form:"-"`
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
	var useTime = t.UpdatedAt.Unix() - t.CreatedAt.Unix()
	if useTime < 0 {
		useTime = 0
	}
	return json.Marshal(struct {
		HistoryMarshalJSON
		UseTime int64 `json:"useTime"`
	}{
		HistoryMarshalJSON(t),
		useTime,
	})
}

func (t *History) Validator() error {
	if t.Version == "" {
		t.Version = "latest"
	}
	if t.ProjectId > 0 {
		if i := GetProject(t.ProjectId); i != nil {
			t.Project = i
		}
	}
	if t.Project == nil {
		return fmt.Errorf("部署失败：project id %d not found", t.ProjectId)
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
	if t.IsDockerMode() {
		t.Project.Native.Shell = ""
		t.Project.Native.Volume = VolumeList{}
	} else {
		t.Project.Docker.Shell = ""
		t.Project.Docker.Image = ""
		t.Project.Docker.RunOptions = ""
		t.Project.Docker.Dockerfile = ""
		t.Project.Docker.Volume = VolumeList{}
	}
	return nil
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

func (t *History) WorkspaceVolumePath(name string) string {
	return path.Join(t.Workspace(), DockerVolumePathName, name)
}

func (t *History) WorkspaceFollowLog() string {
	return path.Join(t.Workspace(), ".follow.log")
}

func (t *History) WorkspaceEndLog() (err error) {
	if err = os.MkdirAll(t.Workspace(), os.ModePerm); err != nil {
		return
	}
	if err = ioutil.WriteFile(t.WorkspaceFollowLog(), []byte(t.Log), os.ModePerm); err != nil {
		return
	}
	return
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

func (t *History) IsDockerMode() bool {
	return t.Project.Mode == ProjectModeDocker || t.Project.Mode == ProjectModeImage
}

func (t *History) IsDockerImageMode() bool {
	return t.Project.Mode == ProjectModeImage
}

func (t *History) getBeforeHistory() *History {
	var res History
	if DB().Last(&res, "project_id=? AND id<?", t.ProjectId, t.ID).Error != nil {
		return nil
	}
	return &res
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
		//duration wait check
		var startTime = time.Now()
		for {
			if time.Now().After(startTime.Add(5 * time.Second)) {
				break
			}
			if t.Status == HistoryStatusDefault {
				time.Sleep(time.Millisecond * 100)
			} else {
				break
			}
		}
	} else {
		t.updateStatus(nil)
	}
	return nil
}

//project online
func (t *History) Online(async bool) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	t.onlineCtx = ctx
	t.onlineCtxCancel = cancel
	onlineMaps.Store(t.ID, t)

	//create workspace
	if err = os.MkdirAll(t.WorkspacePath(DockerVolumePathName), os.ModePerm); err != nil {
		return
	}
	//create run.sh content
	var runContent string
	if t.IsDockerMode() { //deploy mode docker
		runContent, err = t.createRunDockerMode()
	} else { //node shell
		runContent, err = t.createRunNativeMode()
	}
	fmt.Println(runContent)
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
	cmd := exec.Command("sh", "-e", t.WorkspaceRun())
	cmd.Stdout = fileStream
	cmd.Stderr = fileStream
	if async { //async for online apply
		go func() {
			defer fileStream.Close()
			t.updateStatus(cmd)
			cancel()
		}()
	} else { //sync for cronjob online
		defer fileStream.Close()
		t.updateStatus(cmd)
		cancel()
	}

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
					log.StdOut("history", "shutdown", t.WorkspaceRun(), core.KillProcessGroup(cmd))
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

func (t *History) createRunDockerMode() (runContent string, err error) {
	//需要检测之前的history如果存在project.name不一致需要先移除container
	if beforeHistory := t.getBeforeHistory(); beforeHistory != nil {
		if beforeHistory.Project.Name != t.Project.Name {
			if err = beforeHistory.Remove(); err != nil {
				return
			}
		}
	}

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
docker login %s --username=%s --password=%s
docker pull %s
docker build --platform=linux/amd64 -t %s %s
docker push %s
`,
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			fromImage,
			imageName, t.Workspace(),
			imageName,
		)
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
				"docker run -i --name %s %s %s",
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			imageName,
			t.Project.Name,
			t.Project.Name, t.Project.Docker.RunOptions, imageName,
		)
		var sshArgs []string
		sshArgs, err = t.Node.RunSshArgs(false, "", fmt.Sprintf("'%s'", dockerRun))
		if err != nil {
			return
		}
		sshDockerRun = strings.Join(sshArgs, " ")
	}

	//create run.sh content
	runContent = fmt.Sprintf(
		"#!/bin/sh\ncd %s\n%s\n%s\n%s\ncd %s\n",
		t.Workspace(),
		t.Project.Docker.Shell,
		dockerBuild,
		sshDockerRun,
		t.Workspace(),
	)
	return
}

func (t *History) createRunNativeMode() (runContent string, err error) {
	var node = t.Node
	if node == nil {
		err = errors.New("node is nil")
		return
	}
	var template = t.Project.Native
	var scpShell string
	var scpRemoteFilePrefix = fmt.Sprintf("/tmp/devops-%d-%d-", t.ProjectId, t.ID)
	for _, v := range template.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		var fileName = uuid.GetUUID(t.ID, "@", v.Path)
		var localPath = t.WorkspaceVolumePath(fileName)
		var remoteNodePath = v.Path
		if err = ioutil.WriteFile(localPath, []byte(volumeContent), os.ModePerm); err != nil {
			return
		}
		var scpArgs []string
		scpArgs, err = node.RunScpArgs(localPath, remoteNodePath)
		if err != nil {
			return
		}
		scpShell += strings.Join(scpArgs, " ") + "\n"
	}
	//init facade content
	var localFacadePath = t.WorkspaceVolumePath("run")
	var remoteFacadePath = fmt.Sprintf("%s%s", scpRemoteFilePrefix, "run")
	var remoteFacadeContent = fmt.Sprintf("#!/bin/sh\n%s\nrm -rf %s*\n", template.Shell, scpRemoteFilePrefix)
	if err = ioutil.WriteFile(localFacadePath, []byte(remoteFacadeContent), os.ModePerm); err != nil {
		return
	}
	scpRunArgs, err := node.RunScpArgs(localFacadePath, remoteFacadePath)
	if err != nil {
		return
	}
	scpShell += strings.Join(scpRunArgs, " ") + "\n"
	sshArgs, err := node.RunSshArgs(false, "", fmt.Sprintf("'sh -ex %s'", remoteFacadePath))
	if err != nil {
		return
	}
	//create run.sh content
	runContent = fmt.Sprintf(
		"#!/bin/sh\ncd %s\n%s\n%s\ncd %s\n",
		t.Workspace(),
		scpShell,
		strings.Join(sshArgs, " "),
		t.Workspace(),
	)
	return
}

//移除上线
func (t *History) Remove(option ...interface{}) error {
	if err := CronjobStop(t.ProjectId); err != nil {
		return err
	}
	if t.Project.Mode != ProjectModeDocker {
		return nil
	}
	if _, err := t.Node.Exec(fmt.Sprintf("docker rm -f %s", t.Project.Name)); err != nil {
		return err
	}
	return DB(option...).Model(t.Project).Where("id=?", t.ProjectId).Update("deleted", 1).Error
}
