package tty

import (
	"app/library/ginutil"
	"app/library/uuid"
	"app/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type TTYCode string

const (
	TTYCodeBash = "bash" //bash
	TTYCodeNode = "node" //node ssh
	TTYCodeLogs = "logs" //docker logs
	TTYCodeExec = "exec" //docker exec
	TTYCodeTail = "tail" //history apply tail
)

type valCreate struct {
	Code TTYCode `form:"code" binding:"required"`
	ID   uint32  `form:"id"`
}

func Create(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valCreate
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var port int
	var md5ID = uuid.GetUUID(val.Code)
	switch val.Code {
	case TTYCodeBash:
		port, err = models.CreateGoTTY(true, "", "bash")
		//default: "--close-signal", "1", // SIGHUP
	case TTYCodeNode: //node
		if node := models.GetNode(val.ID); node != nil {
			var args []string
			args, err = node.RunArgs("")
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(true, md5ID, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeExec: //docker exec
		if h := models.GetHistory(val.ID); h != nil {
			var args []string
			args, err = h.Node.RunArgs(fmt.Sprintf("docker exec -it %s sh", h.Project.Name))
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(true, md5ID, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeLogs: //docker logs
		if h := models.GetHistory(val.ID); h != nil {
			var args []string
			args, err = h.Node.RunArgs(fmt.Sprintf("docker logs -f -n 1000 %s", h.Project.Name))
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(true, md5ID, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeTail: //history apply tail
		if h := models.GetHistory(val.ID); h != nil {
			var logFilePath = h.WorkspaceFollowLog()
			if h.IsEnd() {
				if endLogPath := h.WorkspaceEndLog(); endLogPath == "" {
					return
				} else {
					logFilePath = endLogPath
				}
			}
			port, err = models.CreateGoTTY(
				false, "",
				"--close-signal", "2", // SIGINT, ctrl-c
				"tail", "-f", "-n", "5000", logFilePath,
			)
		}
	default:
		return nil, errors.New("not port code " + string(val.Code))
	}
	res["port"] = port
	return
}
