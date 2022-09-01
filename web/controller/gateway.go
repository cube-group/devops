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

	detailGroup := group.Group("/i/:nid", middleware.Node())
	detailGroup.GET(".", t.info)
	detailGroup.DELETE(".", t.del)
	detailDockerGroup := detailGroup.Group("/docker")
	detailDockerGroup.GET(".", t.dockerIndex)
	detailDockerGroup.POST("/restart", t.dockerRestart)
	detailDockerGroup.POST("/rm", t.dockerRm)
	detailDockerGroup.POST("/ps", t.dockerPs)
	detailDockerGroup.POST("/stats", t.dockerStats)

	dockerGroup := group.Group("/docker")
	dockerGroup.POST("/version", t.dockerVersion)
}

func (t *NodeController) index(c *gin.Context) {
	res := node.List(c)
	g.HTML(c, "node/index.html", res)
}

func (t *NodeController) create(c *gin.Context) {
	g.HTML(c, "node/info.html", gin.H{})
}

func (t *NodeController) dockerIndex(c *gin.Context) {
	g.HTML(c, "node/docker.html", gin.H{
		"node": models.GetNode(c),
	})
}

//node docker ps
func (t *NodeController) dockerPs(c *gin.Context) {
	res, err := node.DockerPs(c)
	ginutil.JsonAuto(c, "Success", err, res)
}

//node docker restart
func (t *NodeController) dockerRestart(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", node.DockerRestart(c), nil)
}

//node docker rm -f
func (t *NodeController) dockerRm(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", node.DockerRm(c), nil)
}

//node docker version
func (t *NodeController) dockerVersion(c *gin.Context) {
	res, err := node.DockerVersion(c)
	ginutil.JsonAuto(c, "Success", err, res)
}

//node docker stats
func (t *NodeController) dockerStats(c *gin.Context) {
	res, err := node.DockerStats(c)
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
