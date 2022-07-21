package models

import (
	"app/library/log"
	"app/library/types/convert"
	"app/library/types/jsonutil"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"gorm.io/gorm"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	TTY_PORT_START = 40000
)

type TtyPort struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Uid       uint32         `gorm:"" json:"-" form:"-"`
	Port      uint32         `gorm:"" json:"-" form:"-"`
	Cmd       string         `gorm:"" json:"-" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *TtyPort) TableName() string {
	return "d_tty_port"
}

func CreateGoTTY(c *gin.Context, writeFlag bool, arg ...string) (port uint32, err error) {
	//gotty port save
	err = DB().Transaction(func(tx *gorm.DB) error {
		var find TtyPort
		if tx.Order("port DESC").First(&find).Error == nil {
			port = find.Port + 1
		} else {
			port = TTY_PORT_START
		}
		//create tty args
		var appendArg = []string{}
		if writeFlag {
			appendArg = append(appendArg, "-w")
		}
		appendArg = append(appendArg, "-p", convert.MustString(port), "--once", "--timeout", "15")
		arg = append(appendArg, arg...)
		//save tty port
		return tx.Save(&TtyPort{Uid: UID(c), Port: port, Cmd: strings.Join(append([]string{"gotty"}, arg...), " ")}).Error
	})
	if err != nil {
		return
	}
	//create gotty process
	var waitChan = make(chan int, 1)
	defer close(waitChan)
	var cmd *exec.Cmd
	go func() {
		defer func() {
			if er := recover(); er != nil {
				log.StdWarning("gotty", jsonutil.ToString(arg), er)
			}
		}()
		cmd = exec.Command("gotty", arg...)
		log.StdOut("gotty", port, "end", cmd.Run())
		time.Sleep(time.Second)                            //wait for gotty process finish
		DB().Unscoped().Delete(&TtyPort{}, "port=?", port) //delete port maps
	}()
	//test connect
	var testing = true
	go func() {
		var requestURL = fmt.Sprintf("http://127.0.0.1:%d", port)
		var startTime = time.Now()
		for {
			if !testing {
				return
			}
			if resp, er := req.Get(requestURL); er == nil && resp.Response().StatusCode == http.StatusOK { //http ok
				waitChan <- 1
				return
			}
			if time.Now().After(startTime.Add(time.Second * 10)) { //timeout
				err = errors.New("wait gotty timeout")
				if cmd != nil && cmd.Process != nil {
					log.StdWarning("gotty", "timeout killed", cmd.Process.Signal(syscall.SIGINT), cmd.Process.Kill())
				}
				waitChan <- 1
				return
			}
			time.Sleep(time.Millisecond)
		}
	}()
	//end
	<-waitChan
	testing = false
	return
}
