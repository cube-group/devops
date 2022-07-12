package setting

import (
	"app/library/consts"
	"app/library/env"
	"app/library/log"
	"embed"
	"os/exec"
)

var SysWebDebug bool
var SysWebPort string
var SysWebPortTls string

var SqlHost string
var SqlPort int
var SqlUsername string
var SqlPassword string
var SqlDatabase string
var SqlPoolMaxOpen int
var SqlPoolMaxIdle int
var SqlDebug int

var SysGoTtyHost = "127.0.0.1"

//var SysGoTtyPortSshpass string
//var SysGoTtyPortBash string
//var SysGoTtyRandBasicAuth = fmt.Sprintf("%d:%d", rand.IntnRange(100, 100000), rand.IntnRange(100, 100000))

//init env params
func Init(fs map[string]embed.FS) {
	initEmbed(fs)
	initLocal()
	initEnv()
	initCmd()
}

func initEnv() {
	SysWebDebug = env.GetInt(consts.WEB_DEBUG, 1) == 1
	SysWebPort = env.GetString(consts.WEB_PORT, "80")
	SysWebPortTls = env.GetString(consts.WEB_PORT_TLS, "443")
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
	if err := exec.Command("bash", "--version").Run(); err != nil {
		log.StdWarning("init", "cmd/bash", err.Error())
	}
	if err := exec.Command("sshpass", "-V").Run(); err != nil {
		log.StdWarning("init", "cmd/sshpass", err.Error())
	}
	if err := exec.Command("ssh", "-V").Run(); err != nil {
		log.StdWarning("init", "cmd/ssh", err.Error())
	}
	if err := exec.Command("gotty", "-v").Run(); err != nil {
		log.StdWarning("init", "cmd/gotty", err.Error())
	}
}
