package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Tag() gin.HandlerFunc {
	return func(c *gin.Context) {
		if id := convert.MustUint32(c.Param("tag_id")); id > 0 {
			if i := models.GetTag(id); i != nil {
				c.Set(consts.ContextTag, i)
				c.Next()
				return
			}
		}

		g.Warning(c, errors.New("tag resource not found"))
	}
}
