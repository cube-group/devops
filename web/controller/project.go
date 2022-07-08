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
	detailGroup.Use(middleware.ProjectPermission())
	detailGroup.GET(".", t.info)
	detailGroup.GET("/apply", t.apply)
	detailGroup.POST("/online", middleware.Block(), t.online)
	detailGroup.POST("/offline", middleware.Block(), t.offline)
	detailGroup.GET("/pod", t.pod)
	detailGroup.GET("/member", t.member)
	detailGroup.Use(middleware.ProjectOwnPermission())
	detailGroup.POST("/member", t.memberSave)
	detailGroup.DELETE("/pod", t.podDel)
	detailGroup.DELETE(".", t.del)
}

func (t *ProjectController) index(c *gin.Context) {
	res := project.List(c)
	res["registryPath"] = models.CfgRegistryPath()
	res["tags"] = models.TagList()
	g.HTML(c, "project/index.html", res)
}

func (t *ProjectController) create(c *gin.Context) {
	g.HTML(c, "project/info.html", gin.H{})
}

func (t *ProjectController) info(c *gin.Context) {
	g.HTML(c, "project/info.html", gin.H{
		"tags":    models.TagList(),
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

func (t *ProjectController) offline(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", project.Offline(c), nil)
}

func (t *ProjectController) member(c *gin.Context) {
	g.HTML(c, "project/member.html", gin.H{
		"project": models.GetProject(c),
		"owned":   project.MemberList(c),
		"all":     models.GetAllUser(),
	})
}

func (t *ProjectController) memberSave(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", project.MemberSave(c), nil)
}

func (t *ProjectController) pod(c *gin.Context) {
	var obj = models.GetProject(c)
	g.HTML(c, "project/pod.html", gin.H{
		"project": obj,
		"history": obj.GetLatestHistory(),
	})
}

func (t *ProjectController) podDel(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", project.PodDel(c), nil)
}
