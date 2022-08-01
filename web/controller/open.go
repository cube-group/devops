package controller

import (
	"app/library/g"
	"app/web/service/open"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OpenController struct {
}

func (t *OpenController) Init(group *gin.RouterGroup) {
	group.GET("/oauth/callback", t.oauthCallback)
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
