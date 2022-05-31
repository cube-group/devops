package middlewares

import (
	"app/library/crypt/md5"
	"app/library/types/convert"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	AUTH_FORWARDED_GID  = "X-Forwarded-Api-Gid"
	AUTH_FORWARDED_TIME = "X-Forwarded-Api-Time"
	AUTH_FORWARDED_MD5  = "X-Forwarded-Api-Md5"
	AUTH_FORWARDED_AK   = "X-Forwarded-Api-Ak"
)

//微服务网关中间件
func Gateway() gin.HandlerFunc {
	return func(c *gin.Context) {
		gid := convert.MustUint(c.GetHeader(AUTH_FORWARDED_GID))
		timestamp := c.GetHeader(AUTH_FORWARDED_TIME)
		md5Str := c.GetHeader(AUTH_FORWARDED_MD5)

		if gid == 0 || timestamp == "" || md5Str == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "distrusted", "code": 98})
			c.Abort()
			return
		}

		if md5.MD5(fmt.Sprintf("gid=%d&t=%s", gid, timestamp)) != md5Str {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "distrusted", "code": 99})
			c.Abort()
			return
		}

		c.Next()
	}
}
