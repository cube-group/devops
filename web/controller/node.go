package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/middleware"
	"app/web/service/node"
	"github.com/gin-gonic/gin"
)

type NodeController struct {
}

func (t *NodeController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.GET("/create", t.create)
	group.POST("/save", t.save)
	group.POST("/docker/version", t.dockerVersion)
	group.POST("/docker/stats", t.dockerStats)
	detailGroup := group.Group("/i/:nid", middleware.Node())
	detailGroup.GET(".", t.info)
	detailGroup.DELETE(".", t.del)
}

func (t *NodeController) index(c *gin.Context) {
	res := node.List(c)
	g.HTML(c, "node/index.html", res)
}

func (t *NodeController) create(c *gin.Context) {
	g.HTML(c, "node/info.html", gin.H{
	})
}

//node docker version
func (t *NodeController) dockerVersion(c *gin.Context) {
	res, err := node.GetDockerVersion(c)
	ginutil.JsonAuto(c, "Success", err, res)
}

//node docker stats
func (t *NodeController) dockerStats(c *gin.Context) {
	res, err := node.GetDockerStats(c)
	ginutil.JsonAuto(c, "Success", err, res)
}

func (t *NodeController) info(c *gin.Context) {
	g.HTML(c, "node/info.html", gin.H{
		"node": models.GetNode(c),
	})
}

func (t *NodeController) save(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", node.Save(c), nil)
}

func (t *NodeController) del(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", node.Del(c), nil)
}
