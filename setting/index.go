package setting

import (
	"app/library/consts"
	"app/library/env"
	"app/library/log"
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"os/exec"
)

var SysWebDebug bool
var SysWebServer string

var SqlHost string
var SqlPort int
var SqlUsername string
var SqlPassword string
var SqlDatabase string
var SqlPoolMaxOpen int
var SqlPoolMaxIdle int
var SqlDebug int

var SysGoTtyHost = "127.0.0.1"
var SysGoTtyPortSshpass string
var SysGoTtyPortBash string
var SysGoTtyRandBasicAuth = fmt.Sprintf("%d:%d", rand.IntnRange(100, 100000), rand.IntnRange(100, 100000))

//init env params
func Init() {
	initLocal()
	initEnv()
	initCmd()
}

func initEnv() {
	SysWebDebug = env.GetInt(consts.WEB_DEBUG, 0) == 1
	SysWebServer = env.GetString(consts.WEB_SERVER, "0.0.0.0:80")
	SqlHost = env.GetString(consts.DB_HOST, "127.0.0.1")
	SqlPort = env.GetInt(consts.DB_PORT, 3306, 1, 65534)
	SqlUsername = env.GetString(consts.DB_USERNAME, "root")
	SqlPassword = env.GetString(consts.DB_PASSWORD, "root")
	SqlDatabase = env.GetString(consts.DB_DATABASE, "devops")
	SqlPoolMaxIdle = env.GetInt(consts.DB_POOL_MAX_IDLE, 50, 50, 100)
	SqlPoolMaxOpen = env.GetInt(consts.DB_POOL_MAX_OPEN, 200, 200, 500)
	SqlDebug = env.GetInt(consts.DB_DEBUG, 0)
}

func initCmd() {
	// InitCmdDependency
	if err := exec.Command("gotty", "-v").Run(); err != nil {
		log.StdWarning("init", "cmd/gotty", err.Error())
	} else {
		if er := exec.Command("bash", "--version").Run(); er == nil {
			SysGoTtyPortBash = "30000"
		} else {
			log.StdWarning("init", "cmd/bash", er.Error())
		}
		if er := exec.Command("sshpass", "-V").Run(); er == nil {
			SysGoTtyPortSshpass = "30002"
		} else {
			log.StdWarning("init", "cmd/sshpass", er.Error())
		}
	}
}
