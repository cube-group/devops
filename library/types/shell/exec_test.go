package shell

import (
	"fmt"
	"testing"
	"time"
)

func TestExecCmd(t *testing.T) {
	stdOut, stdErr, state := ExecCmd("ls", time.Second, "/")
	fmt.Println("stdOut: ", stdOut)
	fmt.Println("stdErr: ", stdErr)
	fmt.Println("state: ", state)
}
