package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CrossDomainConfig struct {
	AllowOrigins            []string //default *
	AllowHeaders            []string //default *
	ExposeHeaders           []string //default *
	AllowMethods            []string //default *
	AllowCredentialsDisable bool     //允许客户端传递校验信息比如 cookie (重要), default false
}

//跨域中间件
func CrossDomain(cfg CrossDomainConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.AllowOrigins != nil {
			c.Header("Access-Control-Allow-Origin", strings.Join(cfg.AllowOrigins, ","))
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		if cfg.AllowHeaders != nil {
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ","))
		} else {
			c.Header("Access-Control-Allow-Headers", "*")
		}

		if cfg.ExposeHeaders != nil {
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ","))
		} else {
			c.Header("Access-Control-Expose-Headers", "*")
		}

		if cfg.AllowMethods != nil {
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ","))
		} else {
			c.Header("Access-Control-Allow-Methods", "*")
		}

		if cfg.AllowCredentialsDisable {
			c.Header("Access-Control-Allow-Credentials", "false")
		} else {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
