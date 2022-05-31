package middleware

import (
	"github.com/gin-gonic/gin"
)

func Https() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
