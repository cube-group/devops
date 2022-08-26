package models

import (
	"app/library/consts"
	"app/library/crypt/md5"
	"app/library/log"
	"app/library/types/jsonutil"
	"app/library/types/times"
	"app/library/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
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
	//attach nodes
	if res != nil {
		if res.Nodes == nil {
			res.Nodes = NodeList{}
		}
		if res.N.ID > 0 {
			res.Nodes = res.Nodes.Contact(res.N)
		}
	}
	return
}

type HistoryList []History

func (t HistoryList) IDs() (res []uint32) {
	for _, v := range t {
		res = append(res, v.ID)
	}
	return
}

type HistoryMarshalJSON History

type History struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"-"`
	Uid       uint32         `gorm:"" json:"uid" form:"-"`
	N         Node           `gorm:"column:node" json:"-" form:"-" binding:"-"`
	Nodes     NodeList       `gorm:"" json:"nodes" form:"nodes"`
	Version   string         `gorm:"" json:"version" form:"version"`
	Rollback  uint32         `gorm:"" json:"rollback" form:"rollback"`
	Desc      string         `gorm:"" json:"desc" form:"desc" binding:"required"`
	Status    HistoryStatus  `gorm:"status" json:"status" form:"-"`
	ProjectId uint32         `gorm:"" json:"projectId" form:"-"`
	Project   *Project       `gorm:"" json:"project" form:"-"`
	Log       string         `gorm:"" json:"-" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	NodeIds   []uint32       `gorm:"-" json:"nodeIds" form:"nodeIds"`
	pipeline  *Pipeline      `gorm:"-" form:"-"`
	stream    *os.File       `gorm:"-" form:"-"`
}

func (t *History) TableName() string {
	return "d_history"
}

func (t History) MarshalJSON() ([]byte, error) {
	var useTime = t.UpdatedAt.Unix() - t.CreatedAt.Unix()
	if useTime < 0 {
		useTime = 0
	}
	if t.Nodes == nil {
		t.Nodes = NodeList{}
	}
	if t.N.ID > 0 {
		t.Nodes = t.Nodes.Contact(t.N)
	}
	t.Nodes.Security() //for security
	return json.Marshal(struct {
		HistoryMarshalJSON
		UseTime     int64  `json:"useTime"`
		UpdatedTime string `json:"updatedTime"`
	}{
		HistoryMarshalJSON(t),
		useTime,
		times.FormatDatetime(t.UpdatedAt),
	})
}

func (t *History) Validator() error {
	if t.Version == "" {
		t.Version = "latest"
	}
	if t.Project == nil && t.ProjectId > 0 {
		if i := GetProject(t.ProjectId); i != nil {
			t.Project = i
		}
	}
	if t.Project == nil {
		return fmt.Errorf("部署失败：project id %d not found", t.ProjectId)
	}
	if t.Nodes == nil {
		t.Nodes = NodeList{}
	}
	if !t.IsDockerRunMode() || t.Project.IsCronjob() {
		if t.Rollback > 0 {
			return errors.New("该项目模式不支持回滚操作")
		}
	}
	if !t.IsDockerImageMode() {
		if nodes, err := GetSomeNodes(t.NodeIds); err == nil {
			t.Nodes = nodes
		} else {
			return err
		}
		if len(t.Nodes) == 0 {
			return errors.New("部署目标机器不能为空")
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

func (t *History) RollbackImageURL() string {
	if t.Project.Docker.Image != "" {
		return t.Project.Docker.Image
	}
	return t.ImageURL()
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

func (t *History) WorkspaceContainerLog(elem ...interface{}) (logPath string, stream *os.File, err error) {
	logPath = t.WorkspacePath(fmt.Sprintf(".container.log.%s", md5.MD5(jsonutil.ToString(elem))))
	if err = os.MkdirAll(t.Workspace(), os.ModePerm); err != nil {
		return
	}
	stream, err = os.Create(logPath)
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

func (t *History) IsDockerRunMode() bool {
	return t.Project.Mode == ProjectModeDocker
}

func (t *History) IsDockerImageMode() bool {
	return t.Project.Mode == ProjectModeImage
}

//是否可以修改项目名称
func (t *History) CanChangeName(changedName string) bool {
	if t.Project.Name != changedName && t.Nodes.IsActive() {
		return false
	}
	return true
}

func (t *History) Shutdown() error {
	if t.Status != HistoryStatusDefault {
		return errors.New("shutdown failed online is finished")
	}
	if i, ok := onlineMaps.Load(t.ID); ok {
		var target = i.(*History)
		if target.pipeline != nil {
			return target.pipeline.Stop()
		}
	} else {
		t.updateStatus(errors.New("shutdown"))
	}
	onlineMaps.Delete(t.ID)
	return nil
}

//project online
func (t *History) Online(async bool) (err error) {
	//create workspace
	if err = os.MkdirAll(t.WorkspacePath(DockerVolumePathName), os.ModePerm); err != nil {
		return
	}
	//create .follow.log
	t.stream, err = os.Create(t.WorkspaceFollowLog())
	var pipeline *Pipeline
	if t.IsDockerMode() { //deploy mode docker
		pipeline, err = t.createDockerModePipeline()
	} else { //node shell
		pipeline, err = t.createNativeModePipeline()
	}
	if err != nil {
		return
	}
	t.pipeline = pipeline
	//create record
	onlineMaps.Store(t.ID, t)
	if async {
		go pipeline.Run(t.stream, func(err error) {
			t.updateStatus(err)
		})
	} else {
		pipeline.Run(t.stream, func(err error) {
			t.updateStatus(err)
		})
	}
	return nil
}

func (t *History) updateStatus(err error) {
	if t.stream != nil {
		if er := t.stream.Close(); er != nil {
			log.StdWarning("history", "online", "close stream", er)
		}
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
	maps := map[string]interface{}{"status": status}
	if err != nil {
		maps["log"] = string(logBytes) + "\r\n" + err.Error()
	} else {
		maps["log"] = string(logBytes)
	}
	if er := DB().Model(t).Where("id=?", t.ID).Updates(maps).Error; er != nil {
		log.StdWarning("history", "updateStatus", er)
	}
	//delete online history tag
	onlineMaps.Delete(t.ID)
}

func (t *History) createDockerModePipeline() (pipeline *Pipeline, err error) {
	var template = t.Project.Docker
	var imageName string
	var dockerBuild string
	if template.IsNil() {
		err = errors.New("Dockerfile & Image is nil")
		return
	}
	pipeline = NewPipeline()
	if template.Shell != "" {
		pipeline.Push(PipelineStep{Cmd: template.Shell})
	}
	//docker build
	if template.IsBuildAndRun() {
		var dockerfile string
		dockerfile, err = template.GetComplexDockerfile(t.Workspace())
		if err != nil {
			return
		}
		if err = ioutil.WriteFile(t.WorkspaceDockerfile(), []byte(dockerfile), os.ModePerm); err != nil {
			return
		}
		imageName = t.ImageURL()
		dockerBuild = fmt.Sprintf("docker login %s --username=%s --password=%s\n docker build --pull --platform=linux/amd64 -t %s %s\n docker push %s",
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			imageName, t.Workspace(), imageName,
		)
		pipeline.Push(PipelineStep{Cmd: dockerBuild})
	} else {
		imageName = template.Image
	}
	//remote docker run
	if t.Project.Mode == ProjectModeDocker {
		var runOptionsStruct DockerOptionsStruct
		runOptionsStruct, err = template.RunOptions.GetStruct()
		if err != nil {
			return
		}
		runOptionsStruct.Name = t.Project.Name
		runOptionsStruct.Pull = "always"
		runOptionsStruct.Volume = append(runOptionsStruct.Volume, fmt.Sprintf("/data/log/devops/%d/%d:/data/log", t.Project.ID, t.ID))
		dockerRun := fmt.Sprintf(
			"docker login %s --username=%s --password=%s\n"+
				"docker rm -f %s >/dev/null 2>&1\n"+
				"docker run -i %s %s",
			_cfg.RegistryHost, _cfg.RegistryUsername, _cfg.RegistryPassword,
			t.Project.Name, runOptionsStruct.String(), imageName,
		)
		for _, node := range t.Nodes {
			pipeline.Push(PipelineStep{Node: &node, Cmd: dockerRun})
			if port, apiPath, er := template.GetHealth(); er == nil {
				pipeline.Push(PipelineStep{Node: &node, Type: PipeTypeHealth, Health: PipelineHealth{Port: port, Path: apiPath}})
			}
		}
	}
	return
}

func (t *History) createNativeModePipeline() (pipeline *Pipeline, err error) {
	pipeline = NewPipeline()
	var template = t.Project.Native
	var scpRemoteFilePrefix = fmt.Sprintf("/tmp/devops-%d-%d-", t.ProjectId, t.ID)
	var localPaths = make([]string, 0)
	var remotePaths = make([]string, 0)
	for _, v := range template.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		localPath := t.WorkspaceVolumePath(uuid.GetUUID(t.ID, "@", v.Path))
		localPaths = append(localPaths, localPath)
		remotePaths = append(remotePaths, v.Path)
		if err = ioutil.WriteFile(localPath, []byte(volumeContent), os.ModePerm); err != nil {
			return
		}
	}
	//init steps
	for _, node := range t.Nodes {
		for k, _ := range localPaths {
			pipeline.Push(PipelineStep{Node: &node, Type: PipeTypeScp, Scp: PipelineScp{Source: localPaths[k], Target: remotePaths[k]}})
		}
		pipeline.Push(PipelineStep{Node: &node, Cmd: fmt.Sprintf("%s\nrm -rf %s*\n", template.Shell, scpRemoteFilePrefix)})
	}
	return
}

//移除上线
func (t *History) Remove(statusUpdateFlag bool, option ...interface{}) error {
	if err := CronjobStop(t.ProjectId); err != nil {
		return err
	}
	if t.Project.Mode != ProjectModeDocker {
		return nil
	}
	var node *Node
	for _, v := range option {
		switch vv := v.(type) {
		case Node:
			node = &vv
		case *Node:
			node = vv
		}
	}
	if node != nil {
		if _, err := node.Exec(fmt.Sprintf("docker rm -f %s", t.Project.Name)); err != nil {
			return err
		}
		t.Nodes.RemovePod(node.ID)
	} else {
		for k, item := range t.Nodes {
			if _, err := item.Exec(fmt.Sprintf("docker rm -f %s", t.Project.Name)); err != nil {
				return err
			}
			t.Nodes[k].Removed = true
		}
	}
	var db = DB(option...)
	if err := db.Model(t).Where("id=?", t.ID).Update("nodes", t.Nodes).Error; err != nil {
		return err
	}
	if statusUpdateFlag {
		return db.Model(t.Project).Where("id=?", t.ProjectId).Update("deleted", 1).Error
	}
	return nil
}
