package util

import (
	"github.com/gin-gonic/gin"
)

func IsAjax(c *gin.Context) bool {
	if c.GetHeader("X-Requested-With") != "" {
		return true
	}
	if c.GetHeader("x-visible-ajax") == "1" {
		return true
	}
	return false
}
