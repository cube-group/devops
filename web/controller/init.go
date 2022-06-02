package controller

import (
	"app/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/project")
	})
	new(LoginController).Init(e.Group("/login"))
	e.Use(middleware.Auth())
	new(CfgController).Init(e.Group("/cfg"))
	new(TtyController).Init(e.Group("/tty"))
	new(NodeController).Init(e.Group("/node"))
	new(ProjectController).Init(e.Group("/project"))
	new(HistoryController).Init(e.Group("/history"))
}
