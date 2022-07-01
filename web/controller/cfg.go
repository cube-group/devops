package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"app/web/middleware"
	"app/web/service/cfg"
	"github.com/gin-gonic/gin"
)

type CfgController struct {
}

func (t *CfgController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.POST(".", t.save)
	group.GET("/tty", t.localTty)
	tagGroup := group.Group("/tag")
	tagGroup.GET(".", t.tagList)
	tagGroup.POST(".", t.tagSave)
	tagGroup.DELETE("/:tag_id", middleware.Tag(), t.tagDel)
}

func (t *CfgController) index(c *gin.Context) {
	res, _ := models.GetCfg()
	g.HTML(c, "cfg/index.html", gin.H{"cfg": res})
}

func (t *CfgController) save(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", cfg.Save(c), nil)
}

func (t *CfgController) localTty(c *gin.Context) {
	g.HTML(c, "cfg/tty.html", gin.H{})
}

func (t *CfgController) tagList(c *gin.Context) {
	g.HTML(c, "cfg/tag.html", gin.H{"list": cfg.TagList(c)})
}

func (t *CfgController) tagSave(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", cfg.TagSave(c), nil)
}

func (t *CfgController) tagDel(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", cfg.TagDel(c), nil)
}
