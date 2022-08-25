package sshtool

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
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

	_sshClient  *ssh.Client
	_sftpClient *sftp.Client
	_sshSession *ssh.Session
}

func NewSSHClient(host, port, username, password, rsaPrivate string) (res *SSHClient, err error) {
	client := &SSHClient{
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
	c, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return
	}
	client._sshClient = c
	return client, nil
}

func (t *SSHClient) Close() error {
	if t._sshSession != nil {
		t._sshSession.Close()
	}
	if t._sshClient != nil {
		return t._sshClient.Close()
	}
	return nil
}

//ssh multi session exec
func (t *SSHClient) ExecMulti(cmdList []string) (res []SSHClientSessionResult) {
	for _, cmd := range cmdList {
		res = append(res, t.Exec(cmd))
	}
	return
}

//ssh session exec
func (t *SSHClient) Exec(cmd string) (res SSHClientSessionResult) {
	session, err := t._sshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	t._sshSession = session
	resultBytes, err := session.CombinedOutput(cmd)
	return SSHClientSessionResult{
		Result: resultBytes,
		Error:  err,
	}
}

//ssh session exec with stdout & stderr
func (t *SSHClient) ExecWithStd(writer io.Writer, cmd string) (err error) {
	session, err := t._sshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	session.Stdout = writer
	session.Stderr = writer
	t._sshSession = session
	if err = session.Start(cmd); err != nil {
		return
	}
	return session.Wait()
}

func (t *SSHClient) ScpUpload(localPath, remotePath string, autoClose bool) (err error) {
	if t._sftpClient == nil {
		sftpClient, err := sftp.NewClient(t._sshClient)
		if err != nil {
			return err
		}
		t._sftpClient = sftpClient
	}
	if autoClose {
		defer t._sftpClient.Close()
	}

	target, err := t._sftpClient.Create(remotePath)
	if err != nil {
		return
	}
	defer target.Close()

	fileBytes, err := ioutil.ReadFile(localPath)
	if err != nil {
		return
	}

	_, err = target.Write(fileBytes)
	return
}

func (t *SSHClient) ScpDownload(localPath, remotePath string, autoClose bool) (err error) {
	if t._sftpClient == nil {
		sftpClient, err := sftp.NewClient(t._sshClient)
		if err != nil {
			return err
		}
		t._sftpClient = sftpClient
	}
	if autoClose {
		defer t._sftpClient.Close()
	}

	source, err := t._sftpClient.Open(remotePath)
	if err != nil {
		return
	}
	defer source.Close()

	target, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return
}
