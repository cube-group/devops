package models

import (
	"app/library/log"
	"app/library/types/convert"
	"app/library/types/jsonutil"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

var ttyPorts sync.Map

func KillPortProcess(port int) error {
	i, ok := ttyPorts.Load(port)
	if !ok {
		return nil
	}
	cmd := exec.Command("sh", "sh/kill.sh", convert.MustString(i))
	bytes, err := cmd.CombinedOutput()
	fmt.Println(cmd.Args, string(bytes))
	return err
}

func CreateGoTTY(writeFlag bool, md5ID string, arg ...string) (port int, err error) {
	port = 40000
	for {
		if _, ok := ttyPorts.Load(port); ok {
			port++
		} else {
			break
		}
	}
	if port > 50000 {
		err = errors.New("TTY No Port Can Used")
		return
	}
	ttyPorts.Store(port, md5ID)
	var appendArg = []string{}
	if writeFlag {
		appendArg = append(appendArg, "-w")
	}
	appendArg = append(appendArg, "-p", convert.MustString(port), "--once", "--timeout", "10")
	arg = append(appendArg, arg...)

	var waitChan = make(chan int, 1)
	defer close(waitChan)

	//create gotty process
	go func() {
		defer func() {
			if er := recover(); er != nil {
				log.StdWarning("gotty", jsonutil.ToString(arg), er)
			}
		}()
		cmd := exec.Command("gotty", arg...)
		log.StdOut("gotty", port, "end", cmd.Run())
		ttyPorts.Delete(port) //delete port maps
	}()
	//test connect
	var tested = true
	go func() {
		var requestURL = fmt.Sprintf("http://127.0.0.1:%d", port)
		var startTime = time.Now()
		for {
			if !tested {
				return
			}
			if resp, er := req.Get(requestURL); er == nil && resp.Response().StatusCode == http.StatusOK { //http ok
				waitChan <- 1
				return
			}
			if time.Now().After(startTime.Add(time.Second * 3)) { //timeout
				err = errors.New("timeout")
				waitChan <- 1
				return
			}
			time.Sleep(time.Millisecond)
		}
	}()
	<-waitChan
	tested = false
	return
}
