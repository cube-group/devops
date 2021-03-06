package tty

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type TTYCode string

const (
	TTYCodeBash              = "bash"          //bash
	TTYCodeNode              = "node"          //node ssh
	TTYCodeContainerLogs     = "containerLogs" //project node container stdout
	TTYCodeContainerExec     = "containerExec" //project node container exec
	TTYCodeHistoryTail       = "historyTail"   //history apply tail
	TTYCodeNodeContaienrExec = "nodeContainerExec"
	TTYCodeNodeContaienrLogs  = "nodeContainerLogs"
)

type valCreate struct {
	Code TTYCode `form:"code" binding:"required"`
	ID   uint32  `form:"id" binding:"omitempty"`
	Pod  string  `form:"pod" binding:"omitempty"` //for pod exec/log
}

//create gotty process
func Create(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valCreate
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var port uint32
	var args []string
	switch val.Code {
	case TTYCodeContainerExec: //docker container exec
		if h := models.GetHistory(val.ID); h != nil {
			if models.GetUser(c).HasPermissionProject(h.ProjectId) != nil {
				err = errors.New("没有权限操作")
				return
			}
			if node, ok := h.Nodes.Get(convert.MustUint32(val.Pod)); ok {
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
	case TTYCodeContainerLogs: //docker container logs -f
		if h := models.GetHistory(val.ID); h != nil {
			if models.GetUser(c).HasPermissionProject(h.ProjectId) != nil {
				err = errors.New("没有权限操作")
				return
			}
			if node, ok := h.Nodes.Get(convert.MustUint32(val.Pod)); ok {
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
	case TTYCodeHistoryTail: //history apply tail
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
	case TTYCodeBash: //system bash
		if !models.IsAdm(c) {
			err = errors.New("只有管理员可操作性")
			return
		}
		port, err = models.CreateGoTTY(c, true, "bash")
		//default: "--close-signal", "1", // SIGHUP
	case TTYCodeNode: //node ssh
		if !models.IsAdm(c) {
			err = errors.New("只有管理员可操作性")
			return
		}
		if node := models.GetNode(val.ID); node != nil {
			args, err = node.RunSshArgs(true, "", "")
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
			//"--close-signal", "9", // SIGKILL, kill -9
		}
	case TTYCodeNodeContaienrExec: //node container exec
		if !models.IsAdm(c) {
			err = errors.New("没有权限操作")
			return
		}
		if n := models.GetNode(val.ID); n != nil {
			args, err = n.RunSshArgs(true, "", fmt.Sprintf("docker exec -it %s sh", val.Pod))
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
		}
	case TTYCodeNodeContaienrLogs: //node container logs -f
		if !models.IsAdm(c) {
			err = errors.New("没有权限操作")
			return
		}
		if n := models.GetNode(val.ID); n != nil {
			args, err = n.RunSshArgs(false, "", fmt.Sprintf("docker logs -f -n 1000 %s", val.Pod))
			if err != nil {
				return
			}
			port, err = models.CreateGoTTY(c, true, append([]string{"--close-cmd", "exit"}, args...)...)
		}
	default:
		return nil, errors.New("not port code " + string(val.Code))
	}
	res["port"] = port
	return
}
