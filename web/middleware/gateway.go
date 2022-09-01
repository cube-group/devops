package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Node() gin.HandlerFunc {
	return func(c *gin.Context) {
		if id := convert.MustUint32(c.Param("nid")); id > 0 {
			if i := models.GetNode(id); i != nil {
				c.Set(consts.ContextNode, i)
				c.Next()
				return
			}
		}

		g.Warning(c, errors.New("node resource not found"))
	}
}
