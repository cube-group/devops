package controller

import (
	"app/library/g"
	"app/library/ginutil"
	"app/web/service/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (t *UserController) Init(group *gin.RouterGroup) {
	group.GET(".", t.index)
	group.POST(".", t.save)
	group.DELETE("/i/:user_id", t.del)
}

func (t *UserController) index(c *gin.Context) {
	res := user.List(c)
	g.HTML(c, "user/index.html", res)
}

func (t *UserController) save(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", user.Save(c), nil)
}

func (t *UserController) del(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", user.Del(c), nil)
}
