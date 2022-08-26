package models

import (
	"app/library/e"
	"app/library/log"
	"app/library/sshtool"
	"app/library/types/convert"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	TTY_PORT_START = 51111
)

type TTY struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Uid       uint32         `gorm:"" json:"-" form:"-"`
	Port      uint32         `gorm:"" json:"-" form:"-"`
	Cmd       string         `gorm:"" json:"-" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *TTY) TableName() string {
	return "d_tty_port"
}

type TTYCache struct {
	SshClient  *sshtool.SSHClient
	Stream     *os.File
	StreamPath string
}

var ttyMaps sync.Map

func TTYCacheClear(port uint32) {
	i, ok := ttyMaps.Load(port)
	if !ok {
		return
	}
	var cache = i.(*TTYCache)
	if cache.SshClient != nil {
		if err := cache.SshClient.Close(); err != nil {
			log.StdWarning("tty", "cache.close.sshClient", err)
		}
	}
	if cache.Stream != nil {
		if err := cache.Stream.Close(); err != nil {
			log.StdWarning("tty", "cache.close.stream", cache.StreamPath, err)
		}
	}
	if err := os.Remove(cache.StreamPath); err != nil {
		log.StdWarning("tty", "cache.close.streamPath", cache.StreamPath, err)
	}
}

func TTYCreate(c *gin.Context, cache *TTYCache, args ...string) (port uint32, err error) {
	//gotty port save
	err = DB().Transaction(func(tx *gorm.DB) error {
		var find TTY
		if tx.Order("port DESC").First(&find).Error == nil {
			port = find.Port + 1
		} else {
			port = TTY_PORT_START
		}
		//create tty args
		var appendArg = []string{"-p", convert.MustString(port), "--once", "--timeout", "15"}
		args = append(appendArg, args...)
		//save tty port
		return tx.Save(&TTY{Uid: UID(c), Port: port, Cmd: strings.Join(append([]string{"gotty"}, args...), " ")}).Error
	})
	//cache
	ttyMaps.Store(port, cache)
	if err != nil {
		return
	}
	//create gotty process
	var waitChan = make(chan int, 1)
	defer close(waitChan)
	var cmd *exec.Cmd
	go e.TryCatch(func() {
		cmd = exec.Command("gotty", args...)
		log.StdOut("gotty", strings.Join(args, " "), "end", cmd.Run())
		time.Sleep(time.Second)                        //wait for gotty process finish
		DB().Unscoped().Delete(&TTY{}, "port=?", port) //delete port maps
	})
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
