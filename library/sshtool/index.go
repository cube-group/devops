package sshtool

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func SSHConnect(username, password, host, port string) (session *ssh.Session, err error) {
	auth := []ssh.AuthMethod{ssh.PublicKeys(),ssh.Password(password)}
	clientConfig := &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: 10 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return
	}
	session, err = client.NewSession()
	if err != nil {
		return
	}
	return
}
