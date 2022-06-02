package controller

import (
	"app/library/crypt/base64"
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
	group.GET("/local", t.page)
	group.GET("/localws", t.localWs)
	group.GET("/node", t.page)
	group.GET("/nodews", t.nodeWs)
	group.GET("/history", t.page)
	group.GET("/historyws", t.historyWs)
	group.GET("/auth_token.js", t.ttyJs)
	group.GET("/js/*.js", t.ttyJs)
}

//代理tty的所有相对js请求
//举例：http://127.0.0.1:8888/tty/ssh?host=39.106.107.76
//页面请求http://127.0.0.1:8888/tty/auth_token.js
//页面请求http://127.0.0.1:8888/tty/js/hterm.js
//页面请求http://127.0.0.1:8888/tty/js/gotty.js
//页面请求ws://127.0.0.1:8888/tty/sshws
func (t *TtyController) ttyJs(c *gin.Context) {
	ttyUrl := fmt.Sprintf(
		"http://127.0.0.1:%s%s",
		setting.SysGoTtyPortSshpass,
		strings.Split(c.Request.RequestURI, "/tty")[1],
	)
	resp, err := req.Get(ttyUrl, req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, resp.String())
	}
}

func (t *TtyController) page(c *gin.Context) {
	resp, err := req.Get("http://127.0.0.1:"+setting.SysGoTtyPortSshpass, req.Header{"Authorization": "Basic " + base64.Base64Encode(setting.SysGoTtyRandBasicAuth)})
	if err != nil {
		c.String(http.StatusNoContent, "TTY Connect Error "+err.Error())
	} else {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, resp.String())
	}
}

//gotty for remote node ssh
func (t *TtyController) nodeWs(c *gin.Context) {
	tty.Proxy(c.Writer, c.Request, "127.0.0.1:"+setting.SysGoTtyPortSshpass, "sshpass")
}

//gotty for local bash
func (t *TtyController) localWs(c *gin.Context) {
	tty.Proxy(c.Writer, c.Request, "127.0.0.1:"+setting.SysGoTtyPortBash, "local")
}

//gotty for history log
func (t *TtyController) historyWs(c *gin.Context) {
	tty.Proxy(c.Writer, c.Request, "127.0.0.1:"+setting.SysGoTtyPortBash, "history")
}
