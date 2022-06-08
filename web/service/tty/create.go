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
	TTYCodeBash = "bash"
	TTYCodeNode = "node"
	TTYCodeLogs = "logs"
	TTYCodeExec = "exec"
	TTYCodeTail = "tail"
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
		port, err = models.CreateGoTTY(true, md5ID, "bash")
	case TTYCodeNode:
		if node := models.GetNode(val.ID); node != nil {
			port, err = models.CreateGoTTY(
				true,
				md5ID,
				"sshpass", "-p", node.SshPassword,
				"ssh", "-t", "-o", "StrictHostKeyChecking=no", "-p", node.SshPort, node.SshUsername+"@"+node.IP,
				fmt.Sprintf("MD5=%s;bash", md5ID),
			)
		}
	case TTYCodeExec:
		if h := models.GetHistory(val.ID); h != nil {
			port, err = models.CreateGoTTY(
				true,
				md5ID,
				"sshpass", "-p", h.Node.SshPassword,
				"ssh", "-t", "-o", "StrictHostKeyChecking=no", h.Node.SshUsername+"@"+h.Node.IP,
				fmt.Sprintf("MD5=%s;docker exec -it %s sh", md5ID, h.Project.Name),
			)
		}
	case TTYCodeLogs:
		if h := models.GetHistory(val.ID); h != nil {
			port, err = models.CreateGoTTY(
				false,
				"sshpass", "-p", h.Node.SshPassword,
				"ssh", "-o", "StrictHostKeyChecking=no", h.Node.SshUsername+"@"+h.Node.IP,
				fmt.Sprintf("MD5=%s;docker logs -f -n 1000 %s", md5ID, h.Project.Name),
			)
		}
	default:
		return nil, errors.New("not port code " + string(val.Code))
	}
	res["port"] = port
	return
}