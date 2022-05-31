package ginutil

import "github.com/gin-gonic/gin"

func SearchMap(c *gin.Context)map[string]string{
	//search map
	search := make(map[string]string, 0)
	for k, v := range c.Request.URL.Query() {
		if v == nil || len(v) == 0 {
			continue
		}
		search[k] = v[0]
	}
	return search
}
