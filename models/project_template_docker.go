package models

import (
	"app/library/crypt/md5"
	"app/library/types/convert"
	"bufio"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type ProjectTemplateDockerMarshalJSON ProjectTemplateDocker

//k8s project cfg about spec template
type ProjectTemplateDocker struct {
	Shell      string        `json:"shell"`
	Image      string        `json:"image"`
	Dockerfile string        `json:"dockerfile"`
	Health     DockerHealth  `json:"health"`
	Volume     VolumeList    `json:"volume"`
	RandomPort uint32        `json:"randomPort"`
	RunOptions DockerOptions `json:"runOptions"`
}

func (t *ProjectTemplateDocker) Validator() error {
	if s, err := t.RunOptions.Validator(); err != nil {
		return err
	} else {
		if t.IsRandomPort() { //override random ports
			s.Port, _ = s.getRandomPorts()
		}
		t.RunOptions = DockerOptions(s.String())
	}
	for i := 0; i < len(t.Volume); {
		if t.Volume[i].Validator() != nil {
			t.Volume = append(t.Volume[:i], t.Volume[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}

//override marshal json
func (t ProjectTemplateDocker) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProjectTemplateDockerMarshalJSON
	}{
		ProjectTemplateDockerMarshalJSON(t),
	})
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *ProjectTemplateDocker) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t ProjectTemplateDocker) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *ProjectTemplateDocker) IsNil() bool {
	return t.Dockerfile == "" && t.Image == ""
}

func (t *ProjectTemplateDocker) IsBuildAndRun() bool {
	return t.Dockerfile != "" && t.Image == ""
}

func (t *ProjectTemplateDocker) IsHealthCheck() bool {
	return t.Health.Path != "" && t.Health.Port > 0
}

func (t *ProjectTemplateDocker) IsRandomPort() bool {
	return t.RandomPort > 0
}

//获取端口和健康监测信息
func (t *ProjectTemplateDocker) GetPortAndHealth(node Node) (healthURL string, portMapping map[int]int, s DockerOptionsStruct, err error) {
	s, err = t.RunOptions.GetStruct()
	if err != nil {
		return
	}
	var randomNodePorts = make([]int, 0)
	if t.IsRandomPort() {
		if s.Port == nil || len(s.Port) == 0 {
			err = errors.New("随机端口模式下，无法找到容器映射端口配置！")
			return
		}
		randomNodePorts, err = node.GetContainerRandomPort(len(s.Port))
		if err != nil {
			return
		}
	}
	portList, portMapping := s.getRandomPorts(randomNodePorts...)
	//非随机端口，不合法端口检测
	if !t.IsRandomPort() {
		for containerPort, nodePort := range portMapping {
			if containerPort == 0 || nodePort == 0 {
				err = errors.New("非随机端口模式下，宿主机或容器端口需合法")
				return
			}
		}
	}
	s.Port = portList
	if t.IsHealthCheck() {
		if nodePort, ok := portMapping[t.Health.Port]; ok {
			healthURL = fmt.Sprintf("http:127.0.0.1:%d%s", nodePort, t.Health.Path)
		} else {
			err = errors.New("无法匹配健康监测端口，请查看项目端口映射配置")
		}
	}
	return
}

func (t *ProjectTemplateDocker) GetComplexDockerfile(workspace string) (content string, err error) {
	if t.Dockerfile == "" {
		return
	}
	//check dockerfile COPY
	var volumeLines = make([]string, 0)
	for k, v := range t.Volume {
		var volumeContent string
		volumeContent, err = v.Load()
		if err != nil {
			return
		}
		var volumeCopyFileName = md5.MD5(fmt.Sprintf("%d@%s", k, v.Path))
		if err = ioutil.WriteFile(path.Join(workspace, volumeCopyFileName), []byte(volumeContent), os.ModePerm); err != nil {
			return
		}
		volumeLines = append(volumeLines, fmt.Sprintf("COPY %s %s", volumeCopyFileName, v.Path))
	}
	//new Dockerfile
	var dockerfileLines = make([]string, 0)
	var scanner = bufio.NewScanner(strings.NewReader(t.Dockerfile))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var line = scanner.Text()
		var lineScanner = bufio.NewScanner(strings.NewReader(line))
		lineScanner.Split(bufio.ScanWords)
		if lineScanner.Scan() {
			var lineWord = lineScanner.Text()
			if lineWord == "FROM" {
				line = fmt.Sprintf("%s\n%s", line, strings.Join(volumeLines, "\n"))
			}
		}
		dockerfileLines = append(dockerfileLines, line)
	}
	content = strings.Join(dockerfileLines, "\n")
	return
}

type DockerHealth struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *DockerHealth) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t DockerHealth) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type DockerOptionsStruct struct {
	Pull            string   `short:"" long:"pull" description:"Pull image before running (\"always\"|\"missing\"|\"never\") (default \"missing\")"`
	Name            string   `short:"" long:"name" description:"Assign a name to the container"`
	Detach          bool     `short:"d" long:"detach" description:"Run container in background and print container ID"`
	Restart         string   `short:"" long:"restart" description:"restart"`
	Hostname        string   `short:"h" long:"hostname" description:"hostname"`
	Network         string   `short:"" long:"network" description:"Connect a container to a network"`
	Rm              bool     `short:"" long:"rm" description:"Automatically remove the container when it exits"`
	CpuShares       int      `short:"c" long:"cpu-shares" description:"CPU shares (relative weight)"`
	Cpus            int      `short:"" long:"cpus" description:"Number of CPUs"`
	Memory          string   `short:"m" long:"memory" description:" Memory limit"`
	OomKillDisabled bool     `short:"" long:"oom-kill-disabled" description:"Disable OOM Killer"`
	Dns             []string `short:"" long:"dns" description:"Set custom DNS servers"`
	Host            []string `short:"" long:"add-host" description:"Add a custom host-to-IP mapping (host:ip)"`
	Port            []string `short:"p" long:"publish" description:"port"`
	Volume          []string `short:"v" long:"volume" description:"volume"`
	Env             []string `short:"e" long:"env" description:"env"`
	Label           []string `short:"l" long:"label" description:"Set meta data on a container"`
	Link            []string `long:"link" description:"Add link to another container"`
	BlockIoWeight   uint16   `short:"" long:"blkio-weight" description:"Block IO (relative weight), between 10 and 1000, or 0 to disable (default 0)"`
	DeviceWriteBps  []string `short:"" long:"device-write-bps" description:"Limit write rate (bytes per second) to a device (default [])"`
	DeviceReadBps   []string `short:"" long:"device-read-bps" description:"Limit write rate (IO per second) to a device (default [])"`
	DeviceWriteIoPs []string `short:"" long:"device-write-iops" description:"Disable OOM Killer"`
	DeviceReadIoPs  []string `short:"" long:"device-read-iops" description:"Disable OOM Killer"`
}

func (t *DockerOptionsStruct) String() (res string) {
	resSlice := make([]string, 0)
	if t.Detach {
		resSlice = append(resSlice, "-d")
	}
	if t.Restart != "" {
		resSlice = append(resSlice, "--restart="+t.Restart)
	}
	if t.Hostname != "" {
		resSlice = append(resSlice, "--hostname="+t.Hostname)
	}
	if t.Network != "" {
		resSlice = append(resSlice, "--network="+t.Network)
	}
	if t.Rm {
		resSlice = append(resSlice, "--rm")
	}
	if t.CpuShares > 0 {
		resSlice = append(resSlice, fmt.Sprintf("-c %d", t.CpuShares))
	}
	if t.Cpus > 0 {
		resSlice = append(resSlice, fmt.Sprintf("--cpus=%d", t.Cpus))
	}
	if t.Memory != "" {
		resSlice = append(resSlice, fmt.Sprintf("-m %s", t.Memory))
	}
	if t.OomKillDisabled {
		resSlice = append(resSlice, "--oom-kill-disable")
	}
	if t.Dns != nil && len(t.Dns) > 0 {
		for _, v := range t.Dns {
			resSlice = append(resSlice, fmt.Sprintf("--dns=%s", v))
		}
	}
	if t.Host != nil && len(t.Host) > 0 {
		for _, v := range t.Host {
			resSlice = append(resSlice, fmt.Sprintf("--add-host=%s", v))
		}
	}
	if t.Port != nil && len(t.Port) > 0 {
		for _, v := range t.Port {
			resSlice = append(resSlice, fmt.Sprintf("-p %s", v))
		}
	}
	if t.Volume != nil && len(t.Volume) > 0 {
		for _, v := range t.Volume {
			resSlice = append(resSlice, fmt.Sprintf("-v %s", v))
		}
	}
	if t.Env != nil && len(t.Env) > 0 {
		for _, v := range t.Env {
			resSlice = append(resSlice, fmt.Sprintf("-e %s", v))
		}
	}
	if t.Label != nil && len(t.Label) > 0 {
		for _, v := range t.Label {
			resSlice = append(resSlice, fmt.Sprintf("-l %s", v))
		}
	}
	if t.Link != nil && len(t.Link) > 0 {
		for _, v := range t.Link {
			resSlice = append(resSlice, fmt.Sprintf("--link %s", v))
		}
	}
	if t.BlockIoWeight > 0 {
		resSlice = append(resSlice, fmt.Sprintf("--blkio-weight=%d", t.BlockIoWeight))
	}
	if t.DeviceWriteBps != nil && len(t.DeviceWriteBps) > 0 {
		for _, v := range t.DeviceWriteBps {
			resSlice = append(resSlice, fmt.Sprintf("--device-write-bps %s", v))
		}
	}
	if t.DeviceReadBps != nil && len(t.DeviceReadBps) > 0 {
		for _, v := range t.DeviceReadBps {
			resSlice = append(resSlice, fmt.Sprintf("--device-read-bps %s", v))
		}
	}
	if t.DeviceWriteIoPs != nil && len(t.DeviceWriteIoPs) > 0 {
		for _, v := range t.DeviceWriteIoPs {
			resSlice = append(resSlice, fmt.Sprintf("--device-write-iops %s", v))
		}
	}
	if t.DeviceReadIoPs != nil && len(t.DeviceReadIoPs) > 0 {
		for _, v := range t.DeviceReadIoPs {
			resSlice = append(resSlice, fmt.Sprintf("--device-read-iops %s", v))
		}
	}
	if t.Pull != "" {
		resSlice = append(resSlice, "--pull "+t.Pull)
	}
	if t.Name != "" {
		resSlice = append(resSlice, "--name "+t.Name)
	}
	return strings.Join(resSlice, " \\\n")
}

//获取随机端口的替代字符串标识
func (t *DockerOptionsStruct) getRandomPorts(assignPorts ...int) (portList []string, portMapping map[int]int) {
	portList = make([]string, 0)
	portMapping = make(map[int]int)
	if t.Port == nil {
		return
	}
	for k, v := range t.Port {
		portArr := strings.Split(strings.Trim(v, " "), ":")
		if len(portArr) != 2 {
			continue
		}
		containerPort := convert.MustInt(portArr[1])
		if containerPort == 0 {
			continue
		}
		if len(assignPorts) > k {
			portMapping[containerPort] = assignPorts[k]
			portList = append(portList, fmt.Sprintf("%d:%d", assignPorts[k], containerPort))
		} else {
			portMapping[containerPort] = 0
			portList = append(portList, fmt.Sprintf("${randomNodePort}:%d", containerPort))
		}
	}
	return
}

type DockerOptions string

func (t DockerOptions) GetStruct() (s DockerOptionsStruct, err error) {
	var str = strings.Join(strings.Fields((string(t))), " ")
	str = strings.ReplaceAll(str, "\\", "")
	if _, err = flags.ParseArgs(&s, strings.Split(str, " ")); err == nil {
		if s.Port == nil {
			s.Port = make([]string, 0)
		}
	}
	return
}

func (t DockerOptions) Validator() (s DockerOptionsStruct, err error) {
	s, err = t.GetStruct()
	if err != nil {
		return
	}
	for _, v := range s.Volume {
		if strings.Contains(v, ":/data/log") {
			err = errors.New("run options不能包含挂载:/data/log其属于系统自动设置")
			return
		}
	}
	s.Name = ""
	return
}
