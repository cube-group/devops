package ginutil

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetQueryParam(c *gin.Context) gin.H {
	res := make(gin.H)
	for k, v := range c.Request.URL.Query() {
		res[k] = v[0]
	}
	return res
}

func GetNoPageQuery(c *gin.Context) string {
	res := make([]string, 0)
	for k, v := range c.Request.URL.Query() {
		if k == "page" || k == "pageSize" {
			continue
		}
		res = append(res, k+"="+v[0])
	}

	resStr := strings.Join(res, "&")
	if resStr != "" {
		resStr = "&" + resStr
	}
	return resStr
}
