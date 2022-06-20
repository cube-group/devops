package core

import (
	"app/library/types/convert"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"syscall"
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
	bytes, err := exec.Command("pstree", "-p", convert.MustString(pid)).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println("==>", string(bytes))
	re := regexp.MustCompile("[0-9]+")
	pidList := re.FindAllString(string(bytes), -1)
	sort.Strings(pidList)
	for i := len(pidList) - 1; i > 0; i-- {
		var pid = convert.MustInt(pidList[i])
		fmt.Println("===>", pid, syscall.Kill(pid, syscall.SIGTERM))
		break
	}

	//pgid, err := syscall.Getpgid(pid)
	//if err != nil {
	//	return err
	//}
	//note the minus sign
	//return syscall.Kill(-pgid, syscall.SIGTERM)
	return nil
}
