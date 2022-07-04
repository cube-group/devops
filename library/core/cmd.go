package core

import (
	"app/library/types/convert"
	"errors"
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
		time.Sleep(100 * time.Millisecond)
	}

	//ps tree and kill latest pid
	bytes, err := exec.Command("pstree", "-p", convert.MustString(pid)).CombinedOutput()
	if err != nil {
		return err
	}
	re := regexp.MustCompile("[0-9]+")
	pidList := re.FindAllString(string(bytes), -1)
	if pidList != nil && len(pidList) > 0 {
		sort.Strings(pidList)
		syscall.Kill(convert.MustInt(pidList[len(pidList)-1]), syscall.SIGTERM)
	}

	//kill group pid
	//pgid, err := syscall.Getpgid(pid)
	//if err != nil {
	//	return err
	//}
	//note the minus sign
	//return syscall.Kill(-pgid, syscall.SIGTERM)
	return nil
}
