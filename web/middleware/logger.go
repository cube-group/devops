// Author: chenqionghe
// Time: 2018-10
// 中间件：记录请求日志，TODO 待完善

package middleware

import (
	"github.com/gin-gonic/gin"
)

//TODO 记录请求信息到日志，便于排查产问题
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
