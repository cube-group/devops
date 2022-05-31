// Author: chenqionghe
// Time: 2018-10
// json编码

package ginutil

import "github.com/gin-gonic/gin"

//获取请求参数
func RequestParams(c *gin.Context) map[string]interface{} {
	return map[string]interface{}{
		"get":  c.Request.URL.Query(),
		"post": c.Request.PostForm,
	}
}
