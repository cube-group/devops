package ginutil

import (
	"github.com/gin-gonic/gin"
	"strings"
)

//获取用户真实ip
func GetClientIp(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	func() {
		defer func() {
			if e := recover(); e != nil {
				//pass
			}
		}()

		if ip != "" {
			ip = strings.Split(ip, ",")[0]
		}
		if strings.Contains(ip, "127.0.0.1") || ip == "" {
			ip = c.GetHeader("X-real-ip")
		}
		if ip == "" {
			ip = "127.0.0.1"
		}
	}()

	return ip
}
