package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Project() gin.HandlerFunc {
	return func(c *gin.Context) {
		if id := convert.MustUint32(c.Param("pid")); id > 0 {
			if i := models.GetProject(id); i != nil {
				c.Set(consts.ContextProject, i)
				c.Next()
				return
			}
		}

		g.Warning(c, errors.New("project resource not found"))
	}
}
