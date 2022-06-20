package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/middleware"
	"app/web/service/history"
	"github.com/gin-gonic/gin"
)

type HistoryController struct {
}

func (t *HistoryController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.POST("/state", t.state)
	detailGroup := group.Group("/i/:historyId", middleware.History())
	detailGroup.GET(".", t.info)
	detailGroup.POST("/shutdown", t.shutdown)
}

func (t *HistoryController) index(c *gin.Context) {
	res := history.List(c)
	g.HTML(c, "history/index.html", res)
}

func (t *HistoryController) state(c *gin.Context) {
	res, err := history.State(c)
	ginutil.JsonAuto(c, "success", err, res)
}

func (t *HistoryController) info(c *gin.Context) {
	g.HTML(c, "history/info.html", gin.H{
		"history": models.GetHistory(c),
	})
}

func (t *HistoryController) shutdown(c *gin.Context) {
	ginutil.JsonAuto(c, "success", models.GetHistory(c).Shutdown(), nil)
}
