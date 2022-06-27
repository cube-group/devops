package models

import (
	"fmt"
	"strings"
)

//运行本地docker
//并且部署至此node上
func RunDocker(node *Node, containerRunPath, localShell, remoteShell string) (run string, volumes map[string]string, err error) {
	volumes = make(map[string]string)
	volumes["/var/run/docker.sock"] = "/var/run/docker.sock" //for docker in docker
	volumes["/root/.ssh"] = "/root/.ssh"                     //for localShell git

	var args []string
	if node != nil {
		var sshPath string
		sshPath, _, err = node.initReadyIdRsa()
		if err != nil {
			return
		}
		volumes[sshPath] = "/root/.ssh2" //for ssh/scp
		dockerRemoteSshIdRsa := "/root/.ssh2/id_rsa"

		args, err = node.RunSshArgs(dockerRemoteSshIdRsa, remoteShell)
		if err != nil {
			return
		}
	}

	run = fmt.Sprintf(
		"#!/bin/sh\ncd %s\n%s\n%s\ncd %s",
		containerRunPath,
		localShell,
		strings.Join(args, " "),
		containerRunPath,
	)
	fmt.Println(run)
	return
}
