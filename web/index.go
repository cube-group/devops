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
)

func Init() {
	initServer()
}

//init gin server
func initServer() {
	engine := gin.New()
	if setting.SysWebDebug {
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
		var sysWebServerTlsAddress = "0.0.0.0:" + setting.SysWebPortTls
		log.StdOut("Init", "HTTPS", sysWebServerTlsAddress)
		log.StdErr("Init", "TLS", engine.RunTLS(sysWebServerTlsAddress, "ssl.pem", "ssl.key"))
	}()
	var sysWebServerAddress = "0.0.0.0:" + setting.SysWebPort
	log.StdOut("Init", "HTTP", sysWebServerAddress)
	log.StdFatal("Init", "HTTP", engine.Run(sysWebServerAddress))
}

//404页面
func handle404(c *gin.Context) {
	if util.IsAjax(c) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "404！您访问的页面未找到！"})
	} else {
		g.HTML(c, "errors/404.html", nil, http.StatusNotFound)
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
