package middleware

import (
	"app/library/g"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Block() gin.HandlerFunc {
	return func(c *gin.Context) {
		//online block
		if block := models.CfgOnlineBlock(); block != "" {
			g.WarningAsPanic(c, errors.New(block))
		}
		c.Next()
	}
}
