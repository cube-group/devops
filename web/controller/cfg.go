package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/service/cfg"
	"github.com/gin-gonic/gin"
)

type CfgController struct {
}

func (t *CfgController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.POST(".", t.save)
	group.GET("/tty", t.localTty)
}

func (t *CfgController) index(c *gin.Context) {
	res, err := models.GetCfg()
	if err != nil {
		res = gin.H{}
	}
	g.HTML(c, "cfg/index.html", gin.H{"cfg": res})
}

func (t *CfgController) save(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", cfg.Save(c), nil)
}

func (t *CfgController) localTty(c *gin.Context) {
	g.HTML(c, "cfg/tty.html", gin.H{})
}
