package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/web/service/open"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OpenController struct {
}

func (t *OpenController) Init(group *gin.RouterGroup) {
	group.POST("/user/list", t.userList)
	group.GET("/oauth/callback", t.oauthCallback)
}

func (t *OpenController) userList(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", nil, open.UserList(c))
}

func (t *OpenController) oauthCallback(c *gin.Context) {
	ref, err := open.OauthCallback(c)
	if err != nil {
		g.HTML(c, "errors/500.html", gin.H{"content": err.Error()})
	} else {
		if ref == "" {
			ref = "/"
		}
		c.Redirect(http.StatusFound, ref)
	}
}
