package middleware

import (
	"app/library/g"
	"app/library/ginutil"
	"app/library/util"
	"github.com/gin-gonic/gin"
)

func forbiddenOutput(c *gin.Context, err error) {
	if util.IsAjax(c) {
		ginutil.JsonAuto(c, "success", err, nil)
	} else {
		g.HTML(c, "errors/exception.html", gin.H{"title": "æéè­Šćđ€", "content": err.Error()})
	}
}
