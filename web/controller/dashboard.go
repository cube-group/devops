package controller

import (
	"app/library/g"
	"app/models"
	"github.com/gin-gonic/gin"
	"runtime"
)

type DashboardController struct {
}

func (t *DashboardController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
}

func (t *DashboardController) index(c *gin.Context) {
	g.HTML(c, "dashboard/index.html", gin.H{
		"info": gin.H{
			"runTime":      models.SystemRunTime(),
			"goVersion":    runtime.Version(),
			"os":           runtime.GOOS,
			"numGoroutine": runtime.NumGoroutine(),
		},
	})
}
