//Date: 2021-02
//Author: linyang
package middlewares

import (
	"github.com/gin-gonic/gin"
)

//gin限流器
func RateLimiter(maxRequests int) gin.HandlerFunc {
	sem := make(chan struct{}, maxRequests)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}
