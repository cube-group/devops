package controller

import (
	"app/library/ginutil"
	"app/web/service/api"
	"app/web/service/user"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
}

func (t *ApiController) Init(group *gin.RouterGroup) {
	group.POST("/user/list", t.userList)
	group.GET("/img/avatar/:img_id", t.avatarImg)

	dockerGroup := group.Group("/docker")
	dockerGroup.POST("/inspect", t.dockerInspect)
}

func (t *ApiController) userList(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", nil, user.GetList(c))
}

func (t *ApiController) avatarImg(c *gin.Context) {
	api.ImgAvatarGet(c)
}

func (t *ApiController) dockerInspect(c *gin.Context) {
	res, err := api.DockerInspect(c)
	ginutil.JsonAuto(c, "Success", err, res)
}
