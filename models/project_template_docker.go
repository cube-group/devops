package models

import (
	"app/library/crypt/md5"
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
	RunOptions DockerOptions `json:"runOptions"`
	Health     DockerHealth  `json:"health"`
	Volume     VolumeList    `gorm:"" json:"volume" form:"volume"`
}

func (t *ProjectTemplateDocker) Validator() error {
	if err := t.RunOptions.Validator(); err != nil {
		return err
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
	URL string `json:"url"`
	Cmd string `json:"cmd"`
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
	Detach          bool     `short:"d" long:"detach" description:"Run container in background and print container ID"`
	Restart         string   `short:"" long:"restart" description:"restart"`
	Hostname        string   `short:"h" long:"hostname" description:"hostname"`
	Network         string   `short:"" long:"network" description:"Connect a container to a network"`
	Rm              bool     `short:"" long:"rm" description:"Automatically remove the container when it exits"`
	CpuShares       int      `short:"c" long:"cpus" description:"CPU shares (relative weight)"`
	Cpus            int      `short:"" long:"cpu-shares" description:"Number of CPUs"`
	Memory          int      `short:"m" long:"memory" description:" Memory limit"`
	OomKillDisabled bool     `short:"" long:"oom-kill-disabled" description:"Disable OOM Killer"`
	BlockIoWeight   uint16   `short:"" long:"blkio-weight" description:"Block IO (relative weight), between 10 and 1000, or 0 to disable (default 0)"`
	DeviceWriteBps  []string `short:"" long:"device-write-bps" description:"Limit write rate (bytes per second) to a device (default [])"`
	DeviceReadBps   []string `short:"" long:"device-read-bps" description:"Limit write rate (IO per second) to a device (default [])"`
	DeviceWriteIoPs []string `short:"" long:"device-write-iops" description:"Disable OOM Killer"`
	DeviceReadIoPs  []string `short:"" long:"device-read-iops" description:"Disable OOM Killer"`
	Dns             []string `short:"" long:"dns" description:"Set custom DNS servers"`
	Host            []string `short:"" long:"add-host" description:"Add a custom host-to-IP mapping (host:ip)"`
	Port            []string `short:"p" long:"publish" description:"port"`
	Volume          []string `short:"v" long:"volume" description:"volume"`
	Env             []string `short:"e" long:"env" description:"env"`
	Label           []string `short:"l" long:"label" description:"Set meta data on a container"`
}

type DockerOptions string

func (t DockerOptions) Validator() (err error) {
	var str = strings.Join(strings.Fields((string(t))), " ")
	str = strings.ReplaceAll(str, "\\", "")
	var structure DockerOptionsStruct
	_, err = flags.ParseArgs(&structure, strings.Split(str, " "))
	if err != nil {
		return
	}
	for _, v := range structure.Volume {
		if strings.Contains(v, ":/data/log") {
			err = errors.New("run options不能包含挂载:/data/log")
			return
		}
	}
	return
}
