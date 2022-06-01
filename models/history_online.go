package models

import (
	"app/library/sshtool"
	"errors"
	"fmt"
	"os"
)

func (t *History) Online() error {
	var node = GetNode(t.NodeId)
	if node == nil {
		return errors.New("node not found")
	}

	shell, err := t.createSshShell()
	if err != nil {
		return err
	}
	//var stdOut,stdErr bytes.Buffer
	conn, err := sshtool.SSHConnect(node.SshUsername, node.SshPassword, node.IP, node.SshPort)
	if err != nil {
		return err
	}
	conn.Stdout = os.Stdout
	conn.Stderr = os.Stdout

	go func() {
		var err = conn.Run(shell)
		defer conn.Close()
		fmt.Println("----End---", err)
	}()
	return nil
}

func (t *History) createCiShell() (res string, err error) {
	var dockerFlag = t.Project.Docker.Dockerfile != ""
	if !dockerFlag {
		return
	}
	var workspace = fmt.Sprintf("/workspace/%d", t.Project.ID)
	res += fmt.Sprintf("mkdir -p %s", workspace)
	for _, v := range t.Project.Volume {
		var result string
		result, err = v.Load()
		if err != nil {
			return
		}
		res += fmt.Sprintf(`echo '%s' > %s;`, result, v.Path)
	}
	return
}

func (t *History) createSshShell() (res string, err error) {
	if t.Project == nil {
		return
	}
	var dockerFlag = t.Project.Docker.Dockerfile != ""
	if !dockerFlag {
		for _, v := range t.Project.Volume {
			var result string
			result, err = v.Load()
			if err != nil {
				return
			}
			res += fmt.Sprintf(`echo '%s' > %s;`, result, v.Path)
		}
	}

	if t.Project.Shell != "" {
		res += t.Project.Shell + ";"
	}
	return
}
