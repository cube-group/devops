// Author: chenqionghe
// Time: 2018-10
// 标准输出
package g

import (
	"app/library/ginutil"
	"app/models"
	"app/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

//显示HTML，加上头部公共信息(如登录用户)
func HTML(c *gin.Context, template string, data map[string]interface{}, code ...int) {
	c.Header("x-server", "devops")
	c.Header("x-server-version", setting.Version)
	if data == nil {
		data = gin.H{}
	}
	data["_u"] = models.SessionUser(c)
	data["_appVersion"] = setting.Version
	ginutil.HTML(c, template, data, code...)
}

//错误模板
func ErrorHTML(c *gin.Context, title string, err error) {
	c.Header("x-server", "devops")
	c.Header("x-server-version", setting.Version)
	HTML(c, "errors/exception.html", map[string]interface{}{
		"_u":      models.SessionUser(c),
		"title":   title,
		"content": err.Error(),
	}, http.StatusOK)
}
