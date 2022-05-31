// Author: chenqionghe
// Time: 2018-10
// 中间件：框架异常捕获，根据请求方式不同输出友好的提示

package middleware

import (
	"app/library/g"
	"app/library/ginutil"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if re := recover(); re != nil {
				var errorContent string
				if ginutil.Input(c, "debug") == "1" {
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					errorContent = string(buf[:n])
				} else {
					errorContent = fmt.Sprintf("%v", re)
				}
				g.WarningAsPanic(c, errors.New(errorContent))
			}
		}()

		c.Next()
	}
}
