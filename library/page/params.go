package page

import (
	"app/library/types/convert"
	"github.com/gin-gonic/gin"
	"net/url"
)

//获取get参数
func GetParams(c *gin.Context) gin.H {
	u, _ := url.Parse(c.Request.URL.String())
	q := u.Query()
	params := gin.H{}
	for key, values := range q {
		params[key] = values[0]
	}
	return params
}

//获取页数
func GetPage(c *gin.Context) int64 {
	return convert.MustInt64(c.DefaultQuery("page", "1"))

}

//获取分页大小
func GetPageSize(c *gin.Context) int64 {
	return convert.MustInt64(c.DefaultQuery("pagesize", "20"))
}
