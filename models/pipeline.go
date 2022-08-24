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

type PipelineCallback func(err error)

type Pipeline struct {
	logWriter io.Writer
	steps     []PipelineStep
	canceled  bool
}

type PipelineStep struct {
	Node        *Node
	Cmd         string //shell exec
	ServicePort int
	ServicePath string //for healthCheck >=200<300
	parent      *Pipeline
	index       int
	_sshClient  *sshtool.SSHClient
	_cmdClient  *exec.Cmd
}

func NewPipeline(logWriter io.Writer, steps []PipelineStep) *Pipeline {
	return &Pipeline{
		logWriter: logWriter,
		steps:     steps,
	}
}

func (t *Pipeline) Run(callback PipelineCallback) {
	if t.steps == nil || len(t.steps) == 0 {
		callback(errors.New("pipeline no step can used"))
	}
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
	t.canceled = true
	for _, v := range t.steps {
		if err := v.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (t *PipelineStep) Stop() error {
	if t._sshClient != nil {
		return t._sshClient.Close()
	}
	if t._cmdClient != nil {
		return core.KillProcessGroup(t._cmdClient)
	}
	return nil
}

func (t *PipelineStep) Run() error {
	t.log(fmt.Sprintf("+ Step %d", t.index))
	//HealthCheck
	if t.ServicePort > 0 && t.ServicePath != "" && t.Node != nil {
		requestURL := fmt.Sprintf("http://127.0.0.1:%d%s", t.ServicePort, t.ServicePath)
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
	} else if t.Cmd != "" {
		if t.Node != nil {
			sshClient, err := t.Node.NewSshClient()
			if err != nil {
				return err
			}
			t._sshClient = sshClient
			defer sshClient.Close()
			return sshClient.StartAndWait(t.parent.logWriter, t.Cmd)
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
	}
	return errors.New("Invalid Pipeline Step")
}

func (t *PipelineStep) log(content string) {
	_, err := t.parent.logWriter.Write([]byte(content + "\r\n"))
	if err != nil {
		log.StdWarning("pipeline", "log", err)
	}
}
