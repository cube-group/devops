package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/middleware"
	"app/web/service/project"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
}

func (t *ProjectController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.GET("/create", t.create)
	group.POST("/save", t.save)
	detailGroup := group.Group("/i/:pid", middleware.Project())
	detailGroup.GET(".", t.info)
	detailGroup.DELETE(".", t.del)
	detailGroup.GET("/apply", t.apply)
	detailGroup.POST("/online", t.online)
	detailGroup.GET("/pod", t.pod)
	detailGroup.GET("/pod/ws", t.podWs)
}

func (t *ProjectController) index(c *gin.Context) {
	res := project.List(c)
	g.HTML(c, "project/index.html", res)
}

func (t *ProjectController) create(c *gin.Context) {
	g.HTML(c, "project/info.html", gin.H{
	})
}

func (t *ProjectController) info(c *gin.Context) {
	g.HTML(c, "project/info.html", gin.H{
		"project": models.GetProject(c),
	})
}

func (t *ProjectController) save(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", project.Save(c), nil)
}

func (t *ProjectController) del(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", project.Del(c), nil)
}

func (t *ProjectController) apply(c *gin.Context) {
	var obj = models.GetProject(c)
	g.HTML(c, "project/apply.html", gin.H{
		"project": obj,
		"history": obj.GetLatestHistory(),
		"nodes":   models.GetNodes(),
	})
}

func (t *ProjectController) online(c *gin.Context) {
	res, err := project.Online(c)
	ginutil.JsonAuto(c, "Success", err, res)
}

func (t *ProjectController) pod(c *gin.Context) {
	var obj = models.GetProject(c)
	g.HTML(c, "project/pod.html", gin.H{
		"project": obj,
		"history": obj.GetLatestHistory(),
	})
}

func (t *ProjectController) podWs(c *gin.Context) {
	project.PodWs(c)
}
