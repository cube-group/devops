package controller

import (
	"app/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})
	new(LoginController).Init(e.Group("/login"))
	new(OpenController).Init(e.Group("/open"))
	e.Use(middleware.Auth())
	new(DashboardController).Init(e.Group("/dashboard"))
	new(TtyController).Init(e.Group("/tty"))
	new(ProjectController).Init(e.Group("/project"))
	new(HistoryController).Init(e.Group("/history"))
	new(ApiController).Init(e.Group("/api"))
	e.Use(middleware.Adm())
	new(NodeController).Init(e.Group("/node"))
	new(CfgController).Init(e.Group("/cfg"))
	new(UserController).Init(e.Group("/user"))
}
