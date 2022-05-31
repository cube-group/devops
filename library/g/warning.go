package g

import (
	"app/library/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Warning(c *gin.Context, err error) {
	if IsAjax(c) {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
	} else {
		HTML(c, "errors/exception.html", gin.H{"title": "System Warning", "content": err.Error()})
	}
	c.Abort()
}

func WarningAsPanic(c *gin.Context, err error) {
	if IsAjax(c) {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error(), "trace": e.TraceString()})
	} else {
		HTML(c, "errors/exception.html", gin.H{"title": "System Warning", "content": err.Error(), "trace": e.TraceString()})
	}
	c.Abort()
}

func IsAjax(c *gin.Context) bool {
	return c.GetHeader("x-corecd-ajax") == "1"
}
