package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"time"
)

//https://godoc.org/golang.org/x/crypto/ssh
import (
	"net"
)

type Cli struct {
	IP   string //IP地址
	User string //用户名
	Port int    //端口号

	Password   string //密码
	SshKeyPath string //密钥文件

	client       *ssh.Client //ssh客户端
	clientConfig *ssh.ClientConfig
}

//密码连接
func NewByPass(ip, user string, pass string, port int) (c *Cli, err error) {
	c = &Cli{
		IP: ip, User: user, Password: pass, Port: port,
	}
	c.clientConfig = &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	if err = c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

//密钥连接
func NewByKey(ip, user string, sshKeyPath string, port int) (*Cli, error) {
	var c *Cli
	var err error
	c = &Cli{
		IP: ip, User: user, SshKeyPath: sshKeyPath, Port: port,
	}
	key, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	c.clientConfig = &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	if err = c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cli) connect() (err error) {
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	c.client, err = ssh.Dial("tcp", addr, c.clientConfig)
	return err
}

//执行shell
//@param dockerfiles shell脚本命令
func (c Cli) Run(shell string) (stdoutStr, stderrStr string, err error) {
	session, err := c.client.NewSession()
	if err != nil {
		return stdoutStr, stderrStr, err
	}
	defer session.Close()
	var stdout io.Reader
	var stderr io.Reader
	if stdout, err = session.StdoutPipe(); err != nil {
		return stdoutStr, stderrStr, err

	}
	if stderr, err = session.StderrPipe(); err != nil {
		return stdoutStr, stderrStr, err

	}
	err = session.Run(shell)
	stdoutBytes, _ := ioutil.ReadAll(stdout)
	stdoutStr = string(stdoutBytes)
	fmt.Println(stdoutStr)
	stderrBytes, _ := ioutil.ReadAll(stderr)
	stderrStr = string(stderrBytes)
	fmt.Println(stderrStr)
	return stdoutStr, stderrStr, err

}

func SendCommand(in io.WriteCloser, cmd string) error {
	if _, err := in.Write([]byte(cmd + "\n")); err != nil {
		return err
	}

	return nil
}
