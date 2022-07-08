package controller

import (
	"app/library/ginutil"
	"app/web/service/open"
	"github.com/gin-gonic/gin"
)

type OpenController struct {
}

func (t *OpenController) Init(group *gin.RouterGroup) {
	group.POST("/user/list", t.userList)
}

func (t *OpenController) userList(c *gin.Context) {
	ginutil.JsonAuto(c, "Success", nil, open.UserList(c))
}
