package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func History() gin.HandlerFunc {
	return func(c *gin.Context) {
		if id := convert.MustUint32(c.Param("historyId")); id > 0 {
			if i := models.GetHistory(id); i != nil {
				c.Set(consts.ContextHistory, i)
				c.Next()
				return
			}
		}

		g.Warning(c, errors.New("history resource not found"))
	}
}
