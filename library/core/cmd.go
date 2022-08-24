package core

import (
	"app/library/types/convert"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"syscall"
	"time"
)

//杀掉一个进程组
func KillProcessGroup(options ...interface{}) error {
	var pid int
	for _, v := range options {
		switch vv := v.(type) {
		case int:
			pid = vv
		case *exec.Cmd:
			vv.Process.Kill()
			if vv.Process != nil {
				pid = vv.Process.Pid
			}
		}
	}
	if pid == 0 {
		return errors.New("pid is nil")
	}
	//pid kill ctrl+c
	for i := 0; i < 2; i++ {
		syscall.Kill(pid, syscall.SIGINT)
		syscall.Kill(pid, syscall.SIGKILL)
		time.Sleep(100 * time.Millisecond)
	}

	//ps tree and kill latest pid
	bytes, err := exec.Command("pstree", "-p", convert.MustString(pid)).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	re := regexp.MustCompile("[0-9]+")
	pidList := re.FindAllString(string(bytes), -1)
	if pidList != nil && len(pidList) > 0 {
		sort.Strings(pidList)
		syscall.Kill(convert.MustInt(pidList[len(pidList)-1]), syscall.SIGTERM)
	}

	return nil
}
