package web

import (
	"app/library/g"
	"app/library/log"
	"app/library/util"
	"app/setting"
	"app/web/controller"
	"app/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

func Init() {
	initTty()
	initServer()
}

//init goTty server
func initTty() {
	//goTty for bash
	if setting.SysGoTtyPortBash != "" {
		go func() {
			log.StdOut("init", "gotty.bash", setting.SysGoTtyPortBash, "pwd", setting.SysGoTtyRandBasicAuth)
			cmd := exec.Command("gotty", "-w", "-p", setting.SysGoTtyPortBash, "-c", setting.SysGoTtyRandBasicAuth,
				"--title-format", "bash", "--permit-arguments", "bash")
			log.StdWarning("gotty", "bash exit", cmd.Run())
		}()
	}
	//goTty for sshpass
	if setting.SysGoTtyPortSshpass != "" {
		go func() {
			log.StdOut("init", "gotty.sshpass", setting.SysGoTtyPortSshpass, "pwd", setting.SysGoTtyRandBasicAuth)
			cmd := exec.Command("gotty", "-w", "-p", setting.SysGoTtyPortSshpass, "-c", setting.SysGoTtyRandBasicAuth,
				"--title-format", "sshpass", "--permit-arguments", "sshpass")
			log.StdWarning("gotty", "sshpass exit", cmd.Run())
		}()
	}
	//goTty for ssh
	if setting.SysGoTtyPortSsh != "" {
		go func() {
			log.StdOut("init", "gotty.ssh", setting.SysGoTtyPortSsh, "pwd", setting.SysGoTtyRandBasicAuth)
			cmd := exec.Command("gotty", "-w", "-p", setting.SysGoTtyPortSsh, "-c", setting.SysGoTtyRandBasicAuth,
				"--title-format", "ssh", "--permit-arguments", "ssh")
			log.StdWarning("gotty", "ssh exit", cmd.Run())
		}()
	}
}

//init gin server
func initServer() {
	log.StdOut("init", "web.server", setting.SysWebServer)
	engine := gin.New()
	if !setting.SysWebDebug {
		engine.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine.NoRoute(handle404)
	engine.NoMethod(handle404)
	engine.Use(middleware.Recovery())
	initViews(engine)       //init views
	controller.Init(engine) //init controllers

	go func() {
		if _, err := os.Stat("ssl.pem"); err != nil {
			return
		}
		if _, err := os.Stat("ssl.key"); err != nil {
			return
		}
		log.StdErr("Init", "TLS", engine.RunTLS("0.0.0.0:443", "ssl.pem", "ssl.key"))
	}()
	log.StdFatal("Init", "HTTP", engine.Run(setting.SysWebServer))
}

//404页面
func handle404(c *gin.Context) {
	if util.IsAjax(c) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "404！您访问的页面未找到！"})
	} else {
		c.HTML(http.StatusNotFound, "errors/404.html", nil)
	}
}

func initViews(c *gin.Engine) {
	//自定义模板方法函数
	c.SetFuncMap(g.ViewFunc())
	//设置静态文件
	c.Static("/public", "./web/public")
	//加载视图模板
	c.LoadHTMLGlob("web/view/**/*")
}
