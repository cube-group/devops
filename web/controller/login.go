// Author: chenqionghe
// Time: 2018-10
// 用户登录退出

package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/service/sys"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

//控制面板
type LoginController struct{}

//路由初始化
func (t *LoginController) Init(group *gin.RouterGroup) {
	group.GET(".", t.Index)
	group.POST(".", t.Login)
	group.GET("/out", t.Logout)
}

func (t *LoginController) Index(c *gin.Context) {
	g.HTML(
		c,
		"login/login.html",
		gin.H{
			"gitlabOAuthURL": models.GitlabOAuthURL(ginutil.Input(c, "ref")),
		},
		http.StatusAccepted,
	)
}

func (t *LoginController) Logout(c *gin.Context) {
	sys.Logout(c)
	c.Redirect(http.StatusFound, "/login")
}

//登录
func (t *LoginController) Login(c *gin.Context) {
	if err := sys.Login(c); err != nil {
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	refUrl, _ := url.Parse(c.Request.Referer())
	ref := refUrl.Query().Get("ref")
	if ref == "" {
		ref = "/"
	}
	c.Redirect(http.StatusFound, ref)
}
