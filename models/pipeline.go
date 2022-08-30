package models

import (
	"app/library/core"
	"app/library/log"
	"app/library/sshtool"
	"app/library/types/convert"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type PipeType string

const (
	//cmd exec local shell command or remote shell command
	PipeTypeCmd PipeType = "cmd"
	//remote shell curl health check
	PipeTypeHealth PipeType = "health"
	//remote scp upload or download
	PipeTypeScp PipeType = "scp"
)

type PipelineCallback func(err error)

type Pipeline struct {
	steps    []*PipelineStep
	canceled bool
}

type PipelineScp struct {
	Source   string
	Target   string
	Download bool
}

type PipelineStep struct {
	Type       PipeType
	Node       *Node
	HealthURL  string
	Cmd        string
	Scp        PipelineScp
	_sshClient *sshtool.SSHClient
	_cmdClient *exec.Cmd
	_canceled  bool
}

func NewPipeline() *Pipeline {
	return &Pipeline{steps: make([]*PipelineStep, 0)}
}

func (t *Pipeline) log(writer io.Writer, content string) {
	if writer == nil {
		return
	}
	_, err := writer.Write([]byte(content + "\r\n"))
	if err != nil {
		log.StdWarning("pipeline", "log", err)
	}
}

func (t *Pipeline) Push(step PipelineStep) {
	t.steps = append(t.steps, &step)
}

func (t *Pipeline) Run(writer io.Writer, callback PipelineCallback) {
	var err error
	if t.steps == nil || len(t.steps) == 0 {
		err = errors.New("pipeline no step can used")
		goto Error
	}
	defer t.Stop()

	for k, v := range t.steps {
		if t.canceled {
			err = errors.New("canceled")
			goto Error
		}
		t.log(writer, fmt.Sprintf("+ Pipeline Step %d %s", k, v.Type))
		switch v.Type {
		case PipeTypeScp:
			err = v.RunScp(writer)
		case PipeTypeHealth:
			t.log(writer, fmt.Sprintf("+ HealthCheck %s %s ...", v.Node.IP, v.HealthURL))
			err = v.RunHealth(writer)
			if err == nil {
				t.log(writer, "+ HealthCheck OK")
			} else {
				t.log(writer, fmt.Sprintf("+ HealthCheck Error %v", err))
			}
		//case PipeTypeCmd:
		default:
			err = v.RunCmd(writer)
		}
		if err != nil {
			goto Error
		}
	}
	t.log(writer, "+ Pipeline Success")
	callback(nil)
	return

Error:
	t.log(writer, fmt.Sprintf("+ Pipeline Error %v", err))
	callback(err)
}

func (t *Pipeline) Stop() error {
	if t.canceled {
		return nil
	}
	t.canceled = true
	for _, v := range t.steps {
		if err := v.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (t *PipelineStep) Stop() error {
	t._canceled = true
	if t._sshClient != nil {
		return t._sshClient.Close()
	}
	if t._cmdClient != nil {
		return core.KillProcessGroup(t._cmdClient)
	}
	t._sshClient = nil
	t._cmdClient = nil
	return nil
}

func (t *PipelineStep) RunHealth(writer io.Writer) error {
	if t.Node == nil {
		return errors.New("Node is nil")
	}
	var startTime = time.Now()
	for {
		if t._canceled {
			return errors.New("HealthCheck: canceled")
		}
		if result, err := t.Node.Exec(fmt.Sprintf("curl -I -m 2 -o /dev/null -s -w %%{http_code} %s", t.HealthURL)); err == nil {
			statusCode := convert.MustInt(string(result))
			if statusCode >= 200 && statusCode < 300 {
				return nil
			}
		}
		if time.Now().After(startTime.Add(time.Minute * 10)) {
			return errors.New("HealthCheck timeout")
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func (t *PipelineStep) RunCmd(writer io.Writer) error {
	if t.Node != nil {
		sshClient, err := t.Node.NewSshClient()
		if err != nil {
			return err
		}
		t._sshClient = sshClient
		defer sshClient.Close()
		return sshClient.ExecWithStd(writer, t.Cmd)
	} else {
		cmd := exec.Command("sh", "-c", t.Cmd)
		cmd.Stderr = writer
		cmd.Stdout = writer
		t._cmdClient = cmd
		if err := cmd.Start(); err != nil {
			return err
		}
		return cmd.Wait()
	}
}

func (t *PipelineStep) RunScp(writer io.Writer) error {
	if t.Node == nil {
		return errors.New("Node is nil")
	}
	sshClient, err := t.Node.NewSshClient()
	if err != nil {
		return err
	}
	if t.Scp.Download {
		return sshClient.ScpDownload(t.Scp.Source, t.Scp.Target, true)
	} else {
		return sshClient.ScpUpload(t.Scp.Source, t.Scp.Target, true)
	}
}
