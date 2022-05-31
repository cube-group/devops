package shell

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"time"
)

//调用本机命令
func ExecCmd(cmdName string, timeout time.Duration, arg ...string) (string, string, *os.ProcessState) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, cmdName, arg...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	cmd.Run()
	cancel()
	return stdOut.String(), stdErr.String(), cmd.ProcessState
}
