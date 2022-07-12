package controller

import (
	"app/library/ginutil"
	"app/web/service/api"
	"app/web/service/open"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
}

func (t *ApiController) Init(group *gin.RouterGroup) {
	group.POST("/user/list", t.userList)
	group.GET("/img/avatar/:img_id", t.avatarImg)
}

func (t *ApiController) userList(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", nil, open.UserList(c))
}

func (t *ApiController) avatarImg(c *gin.Context) {
	api.ImgAvatarGet(c)
}
