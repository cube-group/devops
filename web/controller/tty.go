package controller

import (
	"app/library/ginutil"
	"app/setting"
	"app/web/service/tty"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"strings"
)

type TtyController struct {
}

func (t *TtyController) Init(group *gin.RouterGroup) {
	group.POST("/create", t.create)
	//group.GET("/local", t.page)
	//group.GET("/localws", t.localWs)
	//group.GET("/node", t.page)
	//group.GET("/nodews", t.nodeWs)
	//group.GET("/history", t.page)
	//group.GET("/historyws", t.historyWs)
	//group.GET("/pod", t.page)
	//group.GET("/podws", t.podWs)
	//group.GET("/auth_token.js", t.ttyJs)
	//group.GET("/js/*.js", t.ttyJs)

	var portGroup = group.Group("/port/:port")
	portGroup.GET("/auth_token.js", t.portConnectJs)
	portGroup.GET("/js/*.js", t.portConnectJs)
	portGroup.GET("/connect", t.portConnect)
	portGroup.GET("/connectws", t.portConnectWs)
}

func (t *TtyController) create(c *gin.Context) {
	res, err := tty.Create(c)
	ginutil.JsonAuto(c, "success", err, res)
}

//代理tty的所有相对js请求
//举例：http://127.0.0.1:8888/tty/ssh?host=39.106.107.76
//页面请求http://127.0.0.1:8888/tty/auth_token.js
//页面请求http://127.0.0.1:8888/tty/js/hterm.js
//页面请求http://127.0.0.1:8888/tty/js/gotty.js
//页面请求ws://127.0.0.1:8888/tty/sshws
//func (t *TtyController) ttyJs(c *gin.Context) {
//	ttyUrl := fmt.Sprintf(
//		"http://127.0.0.1:%s%s",
//		setting.SysGoTtyPortSshpass,
//		strings.Split(c.Request.RequestURI, "/tty")[1],
//	)
//	fmt.Println("=>", ttyUrl)
//	resp, err := req.Get(ttyUrl, req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
//	if err != nil {
//		c.String(http.StatusNotFound, "")
//	} else {
//		c.Header("Content-Type", "application/javascript")
//		c.String(http.StatusOK, resp.String())
//	}
//}

//func (t *TtyController) page(c *gin.Context) {
//	resp, err := req.Get("http://127.0.0.1:"+setting.SysGoTtyPortSshpass, req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
//	if err != nil {
//		c.String(http.StatusNoContent, "TTY Connect Error "+err.Error())
//	} else {
//		c.Header("Content-Type", "text/html")
//		c.String(http.StatusOK, resp.String())
//	}
//}

//gotty for remote node ssh
//func (t *TtyController) nodeWs(c *gin.Context) {
//	tty.Proxy(c, setting.SysGoTtyPortSshpass, "node")
//}

//gotty for local bash
//func (t *TtyController) localWs(c *gin.Context) {
//	tty.Proxy(c, setting.SysGoTtyPortBash, "local")
//}

//gotty for history log
//func (t *TtyController) historyWs(c *gin.Context) {
//	tty.Proxy(c, setting.SysGoTtyPortBash, "history")
//}

//gotty for pod
//func (t *TtyController) podWs(c *gin.Context) {
//	tty.Proxy(c, setting.SysGoTtyPortSshpass, "pod")
//}

func (t *TtyController) portConnect(c *gin.Context) {
	ttyUrl := fmt.Sprintf("http://%s:%s", setting.SysGoTtyHost, c.Param("port"))
	resp, err := req.Get(ttyUrl) // req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
	if err != nil {
		c.String(http.StatusNoContent, "TTY Connect Error "+err.Error())
	} else {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, resp.String())
	}
}

func (t *TtyController) portConnectJs(c *gin.Context) {
	port := c.Param("port")
	ttyUrl := fmt.Sprintf(
		"http://%s:%s%s",
		setting.SysGoTtyHost, port,
		strings.Split(c.Request.RequestURI, "/tty/port/"+port)[1],
	)
	resp, err := req.Get(ttyUrl) //req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, resp.String())
	}
}

func (t *TtyController) portConnectWs(c *gin.Context) {
	tty.Proxy(c, c.Param("port"), "port")
}
