package sshtool

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHClientSessionResult struct {
	Result []byte
	Error  error
}

type SSHClient struct {
	username   string
	password   string
	rsaPrivate string
	host       string
	port       string

	_client *ssh.Client
}

func NewSSHClient(host, port, username, password, rsaPrivate string) (res *SSHClient, err error) {
	sshClient := &SSHClient{
		host:       host,
		port:       port,
		username:   username,
		password:   password,
		rsaPrivate: rsaPrivate,
	}
	auth := []ssh.AuthMethod{}
	if rsaPrivate != "" {
		var signer ssh.Signer
		signer, err = ssh.ParsePrivateKey([]byte(rsaPrivate))
		if err != nil {
			return
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if username != "" && password != "" {
		auth = append(auth, ssh.Password(password))
	}
	if len(auth) == 0 {
		err = errors.New("password & rsa is nil")
		return
	}
	clientConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            auth,
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return
	}
	sshClient._client = client
	return sshClient, nil
}

func (t *SSHClient) Close() error {
	if t._client != nil {
		return t._client.Close()
	}
	return nil
}

func (t *SSHClient) ExecMulti(cmdList []string) (res []SSHClientSessionResult) {
	if t._client == nil {
		return
	}
	for _, cmd := range cmdList {
		res = append(res, t.Exec(cmd))
	}
	return
}

func (t *SSHClient) Exec(cmd string) (res SSHClientSessionResult) {
	if t._client == nil {
		res.Error = errors.New("ssh client is nil")
		return
	}
	session, err := t._client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	resultBytes, err := session.CombinedOutput(cmd)
	return SSHClientSessionResult{
		Result: resultBytes,
		Error:  err,
	}
}
