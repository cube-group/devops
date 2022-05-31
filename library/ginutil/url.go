package ginutil

import (
	"github.com/gin-gonic/gin"
	"strings"
)

//获取当前上下文的domain
func GetCurrentDomain(c *gin.Context) string {
	return strings.Split(c.Request.Host, ":")[0]
}

//获取当前上下文的根域
func GetRootDomain(c *gin.Context) string {
	var domain = GetCurrentDomain(c)
	switch domain {
	case "127.0.0.1", "localhost":
		return domain
	default:
		domains := strings.Split(domain, ".")
		return strings.Join(domains[len(domains)-2:], ".")
	}
}
