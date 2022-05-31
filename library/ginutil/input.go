package ginutil

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"html"
	"io/ioutil"
)

//防xss攻击
//从get或post获取
func Input(c *gin.Context, key string, defaultOptions ...string) string {
	defaultValue := ""
	if len(defaultOptions) > 0 {
		defaultValue = defaultOptions[0]
	}
	if c == nil || key == "" {
		return defaultValue
	}
	if r := c.PostForm(key); r != "" {
		return html.EscapeString(r)
	}
	if r := c.Query(key); r != "" {
		return html.EscapeString(r)
	}
	return defaultValue
}

//封装gin ShouldBind
//可进行多次bind解析
func ShouldBind(c *gin.Context, val interface{}) error {
	data, err := c.GetRawData()
	if err != nil {
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	err = c.ShouldBind(val)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return err
}
