package tty

import (
	"app/library/ginutil"
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
	ID   uint32  `form:"id" binding:"omitempty"`
	Pod  uint32  `form:"pod" binding:"omitempty"` //for pod exec/log
}

//create gotty process
func Create(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valCreate
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var port uint32
	switch val.Code {
	case TTYCodeBash:
		port, err = models.CreateGoTTY(c, true, "bash")
		//default: "--close-signal", "1", // SIGHUP
	case TTYCodeNode: //node
		//permission check
		if !models.GetUser(c).IsAdm() {
			err = errors.New("只有管理员可操作性")
			return
		}
		if node := models.GetNode(val.ID); node != nil {
			var args []string
			args, err = node.RunSshArgs(true, "", "")
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeExec: //docker exec
		if h := models.GetHistory(val.ID); h != nil {
			//permission check
			if models.GetUser(c).HasPermissionProject(h.ProjectId) != nil {
				err = errors.New("没有权限操作")
				return
			}
			var args []string
			if node, ok := h.Nodes.Get(val.Pod); ok {
				args, err = node.RunSshArgs(true, "", fmt.Sprintf("docker exec -it %s sh", h.Project.Name))
			} else {
				err = errors.New("pod not found")
			}
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeLogs: //docker logs
		if h := models.GetHistory(val.ID); h != nil {
			//permission check
			if models.GetUser(c).HasPermissionProject(h.ProjectId) != nil {
				err = errors.New("没有权限操作")
				return
			}
			var args []string
			if node, ok := h.Nodes.Get(val.Pod); ok {
				args, err = node.RunSshArgs(false, "", fmt.Sprintf("docker logs -f -n 1000 %s", h.Project.Name))
			} else {
				err = errors.New("pod not found")
			}
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeTail: //history apply tail
		if h := models.GetHistory(val.ID); h != nil {
			var logFilePath = h.WorkspaceFollowLog()
			if h.IsEnd() {
				if err = h.WorkspaceEndLog(); err != nil {
					return
				}
			}
			port, err = models.CreateGoTTY(
				c,
				false,
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
