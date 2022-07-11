package sshtool

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func SSHConnect(username, password, rsaPrivate, host, port string) (session *ssh.Session, err error) {
	auth := []ssh.AuthMethod{}
	if rsaPrivate != "" {
		var signer ssh.Signer
		signer, err = ssh.ParsePrivateKey([]byte(rsaPrivate))
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if username != "" && password != "" {
		auth = append(auth, ssh.Password(password))
	}
	if len(auth) == 0 {
		return nil, errors.New("password & rsa is nil")
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
	session, err = client.NewSession()
	if err != nil {
		return
	}
	return
}
