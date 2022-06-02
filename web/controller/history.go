package controller

import (
	"app/library/g"
	"app/models"
	"app/web/middleware"
	"app/web/service/history"
	"github.com/gin-gonic/gin"
)

type HistoryController struct {
}

func (t *HistoryController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	detailGroup := group.Group("/i/:historyId", middleware.History())
	detailGroup.GET(".", t.info)
}

func (t *HistoryController) index(c *gin.Context) {
	res := history.List(c)
	g.HTML(c, "history/index.html", res)
}

func (t *HistoryController) info(c *gin.Context) {
	g.HTML(c, "history/info.html", gin.H{
		"history": models.GetHistory(c),
	})
}