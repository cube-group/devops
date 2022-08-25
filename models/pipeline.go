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
	PipeTypeCmd    PipeType = "cmd"
	PipeTypeHealth PipeType = "health"
	PipeTypeScp    PipeType = "scp"
)

type PipelineCallback func(err error)

type Pipeline struct {
	logWriter io.Writer
	steps     []*PipelineStep
	canceled  bool
}

type PipelineHealth struct {
	Port int
	Path string
}

type PipelineScp struct {
	Source   string
	Target   string
	Download bool
}

type PipelineStep struct {
	Type       PipeType
	Node       *Node
	Health     PipelineHealth
	Cmd        string
	Scp        PipelineScp
	parent     *Pipeline
	index      int
	_sshClient *sshtool.SSHClient
	_cmdClient *exec.Cmd
}

func NewPipeline(logWriter io.Writer, steps []*PipelineStep) *Pipeline {
	return &Pipeline{
		logWriter: logWriter,
		steps:     steps,
	}
}

func (t *Pipeline) Run(callback PipelineCallback) {
	if t.steps == nil || len(t.steps) == 0 {
		callback(errors.New("pipeline no step can used"))
	}
	defer t.Stop()
	for k, v := range t.steps {
		if t.canceled {
			callback(errors.New("canceled"))
			return
		}
		v.index = k
		v.parent = t
		if err := v.Run(); err != nil {
			callback(err)
			return
		}
	}
	callback(nil)
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
	t.parent = nil
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

func (t *PipelineStep) Run() error {
	t.log(fmt.Sprintf("+ Step %d", t.index))
	switch t.Type {
	case PipeTypeHealth:
		if t.Node == nil {
			return errors.New("Node is nil")
		}
		if t.Health.Path == "" || t.Health.Port == 0 {
			return errors.New("HealthCheck invalid")
		}
		requestURL := fmt.Sprintf("http://127.0.0.1:%d%s", t.Health.Port, t.Health.Path)
		t.log(fmt.Sprintf("HealthCheck => %s %s", t.Node.IP, requestURL))
		var startTime = time.Now()
		for {
			if t.parent.canceled {
				return errors.New("HealthCheck: canceled")
			}
			if result, err := t.Node.Exec(fmt.Sprintf("curl -I -m 2 -o /dev/null -s -w %%{http_code} %s", requestURL)); err == nil {
				statusCode := convert.MustInt(string(result))
				if statusCode >= 200 && statusCode < 300 {
					t.log(fmt.Sprintf("HealthCheck <= %s %s OK", t.Node.IP, requestURL))
					return nil
				}
			}
			if time.Now().After(startTime.Add(time.Minute * 10)) {
				t.log(fmt.Sprintf("HealthCheck <= %s %s Timeout", t.Node.IP, requestURL))
				return errors.New("HealthCheck timeout")
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	case PipeTypeCmd:
		if t.Node != nil {
			sshClient, err := t.Node.NewSshClient()
			if err != nil {
				return err
			}
			t._sshClient = sshClient
			defer sshClient.Close()
			return sshClient.ExecWithStd(t.parent.logWriter, t.Cmd)
		} else {
			cmd := exec.Command("sh", "-c", t.Cmd)
			cmd.Stderr = t.parent.logWriter
			cmd.Stdout = t.parent.logWriter
			t._cmdClient = cmd
			if err := cmd.Start(); err != nil {
				return err
			}
			return cmd.Wait()
		}
	case PipeTypeScp:
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
	return errors.New("Invalid Pipeline Step")
}

func (t *PipelineStep) log(content string) {
	_, err := t.parent.logWriter.Write([]byte(content + "\r\n"))
	if err != nil {
		log.StdWarning("pipeline", "log", err)
	}
}
