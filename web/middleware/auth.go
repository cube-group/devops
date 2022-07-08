// Author: chenqionghe
// Time: 2018-10
// 中间件：检测是否登录、权限控制等
package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/ginutil"
	"app/models"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

//登录验证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断登录
		user := models.SessionUser(c)
		if user == nil {
			if g.IsAjax(c) {
				ginutil.JsonError(c, "登录已失效，请重新登录！", nil, consts.CODE_LOGOUT)
			} else {
				c.Redirect(http.StatusFound, "/login?ref="+url.QueryEscape(c.Request.URL.String()))
			}
			c.Abort()
			return
		}
		//action 权限校验
		//action, err := new(models.Action).FindAccess(c)
		//if err != nil {
		//	g.Warning(c, errors.New("抱歉，您没有权限访问！请联系管理员"))
		//	return
		//}
		//action log
		//if action != nil {
		//	go new(models.ActionLog).InsertLog(c, action)
		//}
		c.Next()
	}
}

func Adm() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.SessionUser(c)
		if !user.IsAdm() {
			g.Warning(c, errors.New("只有系统管理员才有权限"))
			return
		}
		c.Next()
	}
}
